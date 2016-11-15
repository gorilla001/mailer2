package smtp

import (
	"fmt"
	"io"
	"net/smtp"

	"github.com/tinymailer/mailer/types"
)

func simpleSendEmail(e *types.MailEntry) error {
	var auth smtp.Auth

	if err := e.Validate(); err != nil {
		return err
	}

	if smtpAuth := e.Server.Auth(); smtpAuth != nil {
		auth = *smtpAuth
	}

	return smtp.SendMail(
		e.Server.HostAddr(),
		auth,
		e.From,
		[]string{e.To},
		[]byte(e.MailContent()),
	)
}

// SendEmail is exported
func SendEmail(e *types.MailEntry) (err error) {
	var (
		SMTPAction = "unknown"
		helo       = "localhost"
		client     *smtp.Client
		w          io.WriteCloser
	)

	defer func() {
		if err != nil {
			err = fmt.Errorf("SMTPAction:%s Error:%v", SMTPAction, err)
		}
	}()

	// verify
	SMTPAction = "validate"
	if err = e.Validate(); err != nil {
		return
	}

	// dial
	SMTPAction = "dial"
	if client, err = smtp.Dial(e.Server.HostAddr()); err != nil {
		return
	}

	// helo
	SMTPAction = "helo"
	if e.Helo != "" {
		helo = e.Helo
	}
	if err = client.Hello(helo); err != nil {
		return
	}

	// startssl
	if e.Server.UseSSL() {
		if ok, _ := client.Extension("STARTTLS"); ok {
			SMTPAction = "startssl"
			if err = client.StartTLS(e.Server.TLSConfig()); err != nil {
				return client.Close()
			}
		} else {
			// server doesn't support STARTSSL
		}
	}

	// auth if specified
	if auth := e.Server.Auth(); auth != nil {
		SMTPAction = "auth"
		if err = client.Auth(*auth); err != nil {
			return
		}
	}

	// sender / recipientlist
	SMTPAction = "from"
	if err = client.Mail(e.From); err != nil {
		return
	}
	SMTPAction = "to"
	if err = client.Rcpt(e.To); err != nil {
		return
	}

	// send the mail body
	SMTPAction = "data"
	w, err = client.Data()
	if err != nil {
		return
	}
	if _, err = w.Write([]byte(e.MailContent())); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}

	// quit
	SMTPAction = "quit"
	return client.Quit()
}
