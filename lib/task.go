package lib

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/smtp"
	"github.com/tinymailer/mailer/types"
)

var (
	errNoRecipient = errors.New("no recipient addresses to be sent")
	errNoMail      = errors.New("no email contents to be sent")
	errNoServer    = errors.New("no avaliable smtp servers")
)

// AddTask is exported
func AddTask(t types.Task) error {
	return db.DB().Insert(db.CTASK, t)
}

// ListTask is exported
func ListTask() ([]types.Task, error) {
	ret := make([]types.Task, 0)
	err := db.DB().All(db.CTASK, nil, &ret)
	return ret, err
}

// GetTask is exported
func GetTask(id bson.ObjectId) (types.Task, error) {
	var ret types.Task
	err := db.DB().One(db.CTASK, db.BSONIDQuery(id), &ret)
	return ret, err
}

// DelTask is exported
func DelTask(id bson.ObjectId) error {
	err := db.DB().RemoveAll(db.CTASK, db.BSONIDQuery(id))
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}

// GetTaskWrapper is exported
func GetTaskWrapper(t types.Task) types.TaskWrapper {
	rec, _ := GetRecipient(t.Recipient)
	svrNames := make([]string, 0, len(t.Servers))
	for _, svrID := range t.Servers {
		svr, _ := GetServer(svrID)
		svrNames = append(svrNames, svr.Name())
	}
	tw := types.TaskWrapper{
		Task:          t,
		RecipientName: rec.Name,
		ServerNames:   svrNames,
	}
	return tw
}

// RunTask is exported
func RunTask(task *types.Task, policy *types.SendPolicy) io.ReadCloser {
	r, w := io.Pipe()

	go func(w io.WriteCloser) {
		defer w.Close()

		var (
			err     error
			sender  = json.NewEncoder(w)
			startAt = time.Now()
			counter struct {
				sync.Mutex
				succ int
				fail int
			}
		)
		defer func() {
			msg := types.TaskProgressMsg{
				Finish:  true,
				Succ:    counter.succ,
				Fail:    counter.fail,
				Elapsed: time.Now().Sub(startAt).String(),
			}
			if err != nil {
				msg.Error = err.Error()
			}
			sender.Encode(msg)
		}()

		var (
			mailTos []string
			mails   []*types.Mail
			sg      *serverGroup
			total   int
		)
		// default policy
		if policy == nil {
			policy = types.DefaultSendPolicy
		}
		// prepare
		mailTos, mails, sg, err = taskPrepare(task, policy)
		if err != nil {
			return
		}
		// disable failure retry if we have only one server avaliable
		if sg.size() == 1 {
			policy.MaxRetry = 0
		}

		// nb of max worker tokens for concurrency
		tokens := make(chan struct{}, policy.MaxConcurrency)

		// seed the random source
		rand.Seed(time.Now().UnixNano())

		var wg sync.WaitGroup
		total = len(mails) * len(mailTos)
		wg.Add(total)
		for _, mailTo := range mailTos {
			for _, mail := range mails {
				// take a token to start up a new worker
				// hang blocked if channel already full
				tokens <- struct{}{}

				go func(mailTo string, mail *types.Mail) {
					// randomize sleep [0-1]s to smooth smtp sever load
					// especially while `policy.MaxConcurrency` pretty large
					time.Sleep(time.Millisecond * time.Duration(100*(rand.Int()%10+1)))

					// release token on finished
					defer wg.Done()
					defer func() {
						// wait delay
						time.Sleep(time.Second * time.Duration(policy.WaitDelay))
						// release a token back while the worker finished
						<-tokens
					}()

					var (
						err      error
						server   = sg.next()
						fromUser = server.AuthUser // we use smtp auth user as smtp session <FROM> user
						startAt  = time.Now()
						detail   = map[string]interface{}{
							"mail_to":      mailTo,
							"mail_id":      mail.ID.Hex(),
							"smtp_server":  server.ID.Hex(),
							"start_at":     startAt,
							"server_ranks": sg.stats(),
							"finished":     false,
						}
					)

					// set elapsed, errmsg, counters on finished
					defer func() {
						var nsucc, nfail int
						detail["time_elapsed"] = time.Now().Sub(startAt).String()
						detail["finished"] = true
						counter.Lock()
						if err != nil {
							counter.fail++
							detail["error"] = err.Error()
						} else {
							delete(detail, "error")
							counter.succ++
						}
						nsucc, nfail = counter.succ, counter.fail
						counter.Unlock()
						sender.Encode(types.TaskProgressMsg{
							Detail: detail,
							Succ:   nsucc,
							Fail:   nfail,
						})
					}()

					// now, time to send our pretty!
					sender.Encode(types.TaskProgressMsg{
						Detail: detail,
					})
					mailEntry := types.NewMailEntry(mail, fromUser, mailTo, "", server)
					err = smtp.SendEmail(mailEntry)

					// upgrade current smtp server and return
					if err == nil {
						sg.upgrade(server.ID)
						return
					}

					// downgrade current smtp server
					if err != nil {
						sg.downgrade(server.ID)
					}

					detail["error"] = err.Error()
					sender.Encode(types.TaskProgressMsg{
						Detail: detail,
					})

					// switch server and retry
					// note: we don't upgrade/downgrae on smtp server while retry
					for i := uint(1); i <= policy.MaxRetry; i++ {
						newServer := sg.next()
						mailEntry.SwitchServer(newServer)

						detail["retry_n"] = i
						detail["smtp_server"] = newServer.ID.Hex()
						sender.Encode(types.TaskProgressMsg{
							Detail: detail,
						})

						err = smtp.SendEmail(mailEntry)
						if err == nil {
							break
						}
					}

				}(mailTo, mail)
			}
		}
		wg.Wait()

	}(w)

	return r
}

func taskPrepare(task *types.Task, policy *types.SendPolicy) ([]string, []*types.Mail, *serverGroup, error) {
	rec, err := GetRecipient(task.Recipient)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(rec.Emails) == 0 {
		return nil, nil, nil, errNoRecipient
	}

	mails := make([]*types.Mail, 0)
	for _, mid := range task.Mails {
		mail, err := GetMail(mid)
		if err != nil {
			return nil, nil, nil, err
		}
		mails = append(mails, &mail)
	}
	if len(mails) == 0 {
		return nil, nil, nil, errNoMail
	}

	servers := make([]*types.SMTPServer, 0)
	for _, sid := range task.Servers {
		server, err := GetServer(sid)
		if err != nil {
			return nil, nil, nil, err
		}
		servers = append(servers, &server)
	}

	sg := newServerGroup(servers)
	if sg.empty() {
		return nil, nil, nil, errNoServer
	}

	return rec.Emails, mails, sg, nil
}
