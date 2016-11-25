package cli

import (
	"os"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
	"github.com/tinymailer/mailer/types"
)

// Load is exported
func Load(typ string, c *cli.Context) (err error) {

	defer func() {
		if err == nil {
			os.Stdout.WriteString("+OK\r\n")
		}
	}()

	switch typ {

	case "server":
		server := types.SMTPServer{
			ID:       bson.NewObjectId(),
			Host:     c.String("host"),
			Port:     c.String("port"),
			AuthUser: c.String("user"),
			AuthPass: c.String("password"),
		}
		return loadServer(server)

	case "recipient":
		emails, err := contentFromFileOrCLI(c.String("emails"))
		if err != nil {
			return err
		}
		emails = strings.Replace(emails, "\n", ",", -1)
		recipient := types.Recipient{
			ID:     bson.NewObjectId(),
			Name:   c.String("name"),
			Emails: strings.Split(emails, ","),
		}
		return loadRecipient(recipient)

	case "mail":
		mailbody, err := contentFromFileOrCLI(c.String("body"))
		if err != nil {
			return err
		}
		mail := types.Mail{
			ID:       bson.NewObjectId(),
			FromName: c.String("from-name"),
			Subject:  c.String("subject"),
			Body:     mailbody,
		}
		return loadMail(mail)
	}

	return nil
}

func loadServer(s types.SMTPServer) error {
	if err := s.Validate(); err != nil {
		return err
	}
	return lib.AddServer(s)
}

func loadRecipient(r types.Recipient) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return lib.AddRecipient(r)
}

func loadMail(m types.Mail) error {
	if err := m.Validate(); err != nil {
		return err
	}
	return lib.AddMail(m)
}
