package cli

import (
	"github.com/tinymailer/mailer/smtp"
)

// RunTest is exported
func RunTest() {
	smtp.SendMail()
}
