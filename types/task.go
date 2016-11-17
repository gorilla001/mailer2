package types

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

var (
	// DefaultSendPolicy is exported
	DefaultSendPolicy = &SendPolicy{
		MaxRetry:  3,
		MaxSwitch: 3,
	}
)

// Task status defination
const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusStopped = "stopped"
)

// Task is exported
type Task struct {
	ID        bson.ObjectId   `bson:"_id" json:"id"`
	Recipient bson.ObjectId   `bson:"recipient" json:"recipient"`
	Servers   []bson.ObjectId `bson:"servers" json:"servers"`
	Mails     []bson.ObjectId `bson:"mails" json:"mails"`
	Policy    *SendPolicy     `bson:"policy" json:"policy"`
	Status    string          `bson:"status"  json:"status"`
}

// SendPolicy is exported
type SendPolicy struct {
	MaxRetry  int
	MaxSwitch int
}

// NewTask is exported
func NewTask(recipient string, servers, mails []string, policy *SendPolicy) (Task, error) {
	task := Task{
		Status: StatusCreated,
	}

	if !bson.IsObjectIdHex(recipient) {
		return task, fmt.Errorf("invalid recipient id %s", recipient)
	}
	task.Recipient = bson.ObjectIdHex(recipient)

	for _, s := range servers {
		if !bson.IsObjectIdHex(s) {
			return task, fmt.Errorf("invalid server id %s", s)
		}
		task.Servers = append(task.Servers, bson.ObjectIdHex(s))
	}

	for _, m := range mails {
		if !bson.IsObjectIdHex(m) {
			return task, fmt.Errorf("invalid mail id %s", m)
		}
		task.Mails = append(task.Mails, bson.ObjectIdHex(m))
	}

	if policy == nil {
		task.Policy = DefaultSendPolicy
	}

	return task, nil
}

// Validate is exported
func (t Task) Validate() error {
	if len(t.Servers) == 0 {
		return errors.New("servers required")
	}
	if len(t.Mails) == 0 {
		return errors.New("mails required")
	}
	return t.Policy.validate()
}

func (p SendPolicy) validate() error {
	return nil
}

// TaskWrapper is for easily readable
type TaskWrapper struct {
	Task
	RecipientName string
	ServerNames   []string
}

// TaskOptions is exported
type TaskOptions struct {
	MaxConcurrency int
}

// TaskProgressMsg is exported
type TaskProgressMsg struct {
	Detail map[string]interface{} `json:"details,omitempty"`
	Finish bool                   `json:"finish,omitempty"`
	Succ   int                    `json:"succ,omitempty"`
	Fail   int                    `json:"fail,omitempty"`
	Error  string                 `json:"error,omitempty"`
}
