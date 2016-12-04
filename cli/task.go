package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
	"github.com/tinymailer/mailer/types"
)

// Task is exported
func Task(typ string, c *cli.Context) error {
	switch typ {

	case "create":
		task, err := types.NewTask(
			c.String("recipient"),
			strings.Split(c.String("servers"), ","),
			strings.Split(c.String("mails"), ","),
		)
		if err != nil {
			return err
		}
		task.ID = bson.NewObjectId()
		return createTask(task)

	case "run":
		bid, err := bsonID(c)
		if err != nil {
			return err
		}
		return runTask(bid)

	case "show":
		bid, _ := bsonID(c)
		return showTask(bid)

	case "follow":
		bid, err := bsonID(c)
		if err != nil {
			return err
		}
		return followTask(bid)

	case "stop":
		bid, err := bsonID(c)
		if err != nil {
			return err
		}
		return stopTask(bid)

	case "rm":
		bid, err := bsonID(c)
		if err != nil {
			return err
		}
		return rmTask(bid)
	}

	return nil
}

func createTask(t types.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if err := lib.AddTask(t); err != nil {
		return err
	}
	os.Stdout.WriteString("+OK\r\n")
	return nil
}

func runTask(id bson.ObjectId) error {
	task, err := lib.GetTask(id)
	if err != nil {
		return err
	}

	// TODO
	// 1. implement task options
	// 2. dispatch stream to multi handlers
	stream := lib.RunTask(&task, nil)
	defer stream.Close()

	var (
		msg     types.TaskProgressMsg
		decoder = json.NewDecoder(stream)
	)
	for {
		err = decoder.Decode(&msg)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream decode error %v", err)
		}

		pretty(msg)

		if msg.Finish {
			if msg.Error == "" {
				return nil
			}
			return errors.New(msg.Error)
		}

		msg = types.TaskProgressMsg{}
	}
	return nil
}

func showTask(id bson.ObjectId) error {
	if id.Valid() {
		t, err := lib.GetTask(id)
		if err != nil {
			return err
		}
		return pretty(lib.GetTaskWrapper(t))
	}

	ts, err := lib.ListTask()
	if err != nil {
		return err
	}
	tws := make([]types.TaskWrapper, 0)
	for _, t := range ts {
		tws = append(tws, lib.GetTaskWrapper(t))
	}
	return pretty(tws)
}

func followTask(id bson.ObjectId) error {
	return errors.New("not implement yet")
}

func stopTask(id bson.ObjectId) error {
	return errors.New("not implement yet")
}

func rmTask(id bson.ObjectId) error {
	if err := stopTask(id); err != nil {
		return err
	}
	if err := lib.DelTask(id); err != nil {
		return err
	}
	os.Stdout.WriteString("+OK\r\n")
	return nil
}
