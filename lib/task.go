package lib

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/smtp"
	"github.com/tinymailer/mailer/types"
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
func RunTask(task *types.Task, opts *types.TaskOptions) io.ReadCloser {
	r, w := io.Pipe()

	go func(w io.WriteCloser) {
		defer w.Close()

		var (
			err     error
			sender  = json.NewEncoder(w)
			msg     = types.TaskProgressMsg{}
			counter struct {
				sync.Mutex
				succ int
				fail int
			}
		)
		defer func() {
			msg.Finish = true
			msg.Succ = counter.succ
			msg.Fail = counter.fail
			if err != nil {
				msg.Error = err.Error()
			}
			sender.Encode(msg)
		}()

		var (
			mailTos []string
			mails   []*types.Mail
			servers []*types.SMTPServer
		)
		mailTos, mails, servers, err = taskPrepare(task)
		if err != nil {
			return
		}

		if len(servers) == 0 {
			err = errors.New("no avaliable smtp servers")
			return
		}

		var (
			wg     sync.WaitGroup
			detail = map[string]interface{}{}
		)
		for _, mailTo := range mailTos {
			for _, mail := range mails {

				wg.Add(1)
				go func(mail *types.Mail) {
					defer wg.Done()

					detail["to"] = mailTo
					mailEntry := types.NewMailEntry(mail, servers[0].AuthUser, mailTo, "", servers[0])
					counter.Lock()
					if err = smtp.SendEmail(mailEntry); err != nil {
						detail["err"] = err.Error()
						counter.fail++
					} else {
						counter.succ++
					}
					counter.Unlock()
					msg.Detail = detail

					sender.Encode(msg)
					msg = types.TaskProgressMsg{} //reset
				}(mail)

			}
		}
		wg.Wait()

	}(w)

	return r
}

func taskPrepare(task *types.Task) ([]string, []*types.Mail, []*types.SMTPServer, error) {
	rec, err := GetRecipient(task.Recipient)
	if err != nil {
		return nil, nil, nil, err
	}

	mails := make([]*types.Mail, 0)
	for _, mid := range task.Mails {
		mail, err := GetMail(mid)
		if err != nil {
			return nil, nil, nil, err
		}
		mails = append(mails, &mail)
	}

	servers := make([]*types.SMTPServer, 0)
	for _, sid := range task.Servers {
		server, err := GetServer(sid)
		if err != nil {
			return nil, nil, nil, err
		}
		servers = append(servers, &server)
	}

	return rec.Emails, mails, servers, nil
}
