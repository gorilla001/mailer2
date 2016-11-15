package cli

import (
	"encoding/json"
	"os"

	"github.com/codegangsta/cli"

	"github.com/tinymailer/mailer/lib"
)

// Show is exported
func Show(typ string, c *cli.Context) error {
	switch typ {

	case "server":
		return showServer()
	case "recipient":
		return showRecipient()
	case "mail":
		return showMail()
	}

	return nil
}

func pretty(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(append(b, '\n'))
	return err
}

func showServer() error {
	ss, err := lib.ListServer()
	if err != nil {
		return err
	}
	return pretty(ss)
}

func showRecipient() error {
	rs, err := lib.ListRecipient()
	if err != nil {
		return err
	}
	return pretty(rs)
}

func showMail() error {
	ms, err := lib.ListMail()
	if err != nil {
		return err
	}
	return pretty(ms)
}
