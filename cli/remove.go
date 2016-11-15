package cli

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// Remove is exported
func Remove(typ string, c *cli.Context) error {
	id := c.String("id")
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("(%s) not a valid bson id", id)
	}
	bid := bson.ObjectIdHex(id)
	switch typ {

	case "server":
		return removeServer(bid)
	case "recipient":
		return removeRecipient(bid)
	case "mail":
		return removeMail(bid)
	}

	return nil
}

func removeServer(id bson.ObjectId) error {
	err := lib.DelServer(id)
	if err != nil {
		return err
	}
	os.Stdout.WriteString("OK\r\n")
	return nil
}

func removeRecipient(id bson.ObjectId) error {
	err := lib.DelRecipient(id)
	if err != nil {
		return err
	}
	os.Stdout.WriteString("OK\r\n")
	return nil
}

func removeMail(id bson.ObjectId) error {
	err := lib.DelMail(id)
	if err != nil {
		return err
	}
	os.Stdout.WriteString("OK\r\n")
	return nil
}
