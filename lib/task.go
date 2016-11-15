package lib

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
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
func RunTask(task *types.Task, opts *types.TaskOptions) error {
	return nil
}
