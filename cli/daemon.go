package cli

import "github.com/tinymailer/mailer/api"

// Daemon start up daemon to serve API requests
func Daemon(listen, pidfile string) error {
	return api.ListenAndServe(listen)
}
