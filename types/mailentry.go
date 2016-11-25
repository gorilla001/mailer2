package types

import (
	"bytes"
	"encoding/base64"
	"errors"
	"mime"
	"net/mail"

	"gopkg.in/mgo.v2/bson"
)

// Mail is exported
type Mail struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	FromName string        `bson:"from_name" json:"from_name"`
	Subject  string        `bson:"subject" json:"subject"`
	Body     string        `bson:"body" json:"body"`
}

// Validate is exported
func (m *Mail) Validate() error {
	if m.Subject == "" {
		return errors.New("mail subject required")
	}
	if m.Body == "" {
		return errors.New("mail body required")
	}
	return nil
}

// MailEntry is exported
type MailEntry struct {
	*Mail
	From   string      `bson:"from" json:"from"`
	To     string      `bson:"to" json:"to"`
	Helo   string      `bson:"helo" json:"helo"`
	Server *SMTPServer `bson:"server" json:"server"`
}

// Validate is exported
func (e *MailEntry) Validate() error {
	if e.Mail == nil {
		return errors.New("nil mail")
	}
	if err := e.Mail.Validate(); err != nil {
		return err
	}
	if e.Server == nil {
		return errors.New("nil smtp server")
	}
	return e.Server.Validate()
}

// NewMailEntry is exported
func NewMailEntry(mail *Mail, from, to, helo string, server *SMTPServer) *MailEntry {
	return &MailEntry{
		Mail:   mail,
		From:   from,
		To:     to,
		Helo:   helo,
		Server: server,
	}
}

// SwitchServer switch to a new smtp server while retrying failed delivery
func (e *MailEntry) SwitchServer(s *SMTPServer) {
	e.Server = s
}

// MailContent generate the whole mail content to be deliveried
func (e *MailEntry) MailContent() string {
	return e.headerString() + "\r\n" + base64.StdEncoding.EncodeToString([]byte(e.Body))
}

func (e *MailEntry) headerString() string {
	buf := bytes.NewBuffer(nil)
	for k, v := range e.headerMap() {
		buf.Write([]byte(k))
		buf.Write([]byte{':', ' '})
		buf.Write([]byte(v))
		buf.Write([]byte{'\r', '\n'})
	}
	return buf.String()
}

func (e *MailEntry) headerMap() map[string]string {
	addr := mail.Address{e.FromName, e.From}
	return map[string]string{
		"From":                      addr.String(), // -> "name" <xx@yy.zz>
		"To":                        e.To,
		"Reply-To":                  e.From,
		"Subject":                   encodeRFC2047(e.Subject),
		"MIME-Version":              "1.0",
		"Content-Type":              `text/html; charset="utf-8"`,
		"Content-Transfer-Encoding": "base64",
	}
}

// use mail's rfc2047 to encode any string
// See: https://godoc.org/mime#pkg-constants
func encodeRFC2047(s string) string {
	return mime.QEncoding.Encode("utf-8", s)
}
