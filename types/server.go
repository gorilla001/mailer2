package types

import (
	"gopkg.in/mgo.v2/bson"

	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"
)

var (
	// DefaultSMTPServer is exported
	DefaultSMTPServer = &SMTPServer{
		Host:     "smtp.126.com",
		Port:     "25",
		AuthUser: "eyou_uetest@126.com",
		AuthPass: "test123",
	}
)

// SMTPServer is exported
type SMTPServer struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Host     string        `bson:"host" json:"host"`
	Port     string        `bson:"port" json:"port"`
	AuthUser string        `bson:"auth_user" json:"auth_user"`
	AuthPass string        `bson:"auth_pass" json:"auth_pass"`
}

// Validate is exported
func (s *SMTPServer) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("host required")
	}
	if s.Port == "" {
		return fmt.Errorf("port required")
	}
	if s.Port != "25" && s.Port != "465" {
		return fmt.Errorf("port [%s] invalid, must be 25 or 465", s.Port)
	}
	if s.AuthUser == "" || s.AuthPass == "" {
		return fmt.Errorf("user & password required")
	}
	return nil
}

// Name is the SMTPServer uniq key's combination
func (s *SMTPServer) Name() string {
	return s.Host + ":" + s.Port + "-" + s.AuthUser
}

// HostAddr is exported
func (s *SMTPServer) HostAddr() string {
	return s.Host + ":" + s.Port
}

// UseSSL is exported
func (s *SMTPServer) UseSSL() bool {
	return s.Port == "465"
}

// TLSConfig is exported
func (s *SMTPServer) TLSConfig() *tls.Config {
	return &tls.Config{
		ServerName: s.Host,
	}
}

// Ping is exported
func (s *SMTPServer) Ping() (time.Duration, error) {
	return time.Second, nil
}

// Auth is exported
func (s *SMTPServer) Auth() *smtp.Auth {
	if s.AuthUser != "" && s.AuthPass != "" {
		smtpAuth := smtp.PlainAuth("", s.AuthUser, s.AuthPass, s.Host)
		return &smtpAuth
	}
	return nil
}
