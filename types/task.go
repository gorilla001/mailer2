package types

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

var (
	// DefaultSendPolicy is exported
	DefaultSendPolicy = &SendPolicy{
		MaxConcurrency: 3,
		MaxRetry:       0,
		WaitDelay:      10,
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
	Status    string          `bson:"status"  json:"status"`
}

// SendPolicy is exported
type SendPolicy struct {
	MaxConcurrency uint `json:"max_concurrency"`
	MaxRetry       uint `json:"max_retry"`
	WaitDelay      uint `json:"wait_delay"`
}

// NewTask is exported
func NewTask(recipient string, servers, mails []string) (Task, error) {
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
	return nil
}

// TaskWrapper is for easily readable
type TaskWrapper struct {
	Task
	RecipientName string   `json:"recipient_name"`
	ServerNames   []string `json:"server_names"`
}

// TaskProgressMsg is exported
type TaskProgressMsg struct {
	Detail  map[string]interface{} `json:"details,omitempty"` // progress details of per task
	Finish  bool                   `json:"finish,omitempty"`  // final(only) finished
	Succ    int                    `json:"succ,omitempty"`    // final succ counter if finished = true, otherwise halfway
	Fail    int                    `json:"fail,omitempty"`    // final fail counter if finished = true, otherwise halfway
	Error   string                 `json:"error,omitempty"`   // final(only) error message
	Elapsed string                 `json:"elapsed,omitempty"` // final(only) elapsed time
}
