package cli

import (
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// Show is exported
func Show(typ string, c *cli.Context) error {
	bid, _ := bsonID(c)

	switch typ {

	case "server":
		return showServer(bid)
	case "recipient":
		return showRecipient(bid)
	case "mail":
		return showMail(bid)
	}

	return nil
}

func showServer(id bson.ObjectId) error {
	var (
		ss  interface{}
		err error
	)
	if id.Valid() {
		ss, err = lib.GetServer(id)
	} else {
		ss, err = lib.ListServer()
	}
	if err != nil {
		return err
	}
	return pretty(ss)
}

func showRecipient(id bson.ObjectId) error {
	var (
		rs  interface{}
		err error
	)
	if id.Valid() {
		rs, err = lib.GetRecipient(id)
	} else {
		rs, err = lib.ListRecipient()
	}
	if err != nil {
		return err
	}
	return pretty(rs)
}

func showMail(id bson.ObjectId) error {
	var (
		ms  interface{}
		err error
	)
	if id.Valid() {
		ms, err = lib.GetMail(id)
	} else {
		ms, err = lib.ListMail()
	}
	if err != nil {
		return err
	}
	return pretty(ms)
}
