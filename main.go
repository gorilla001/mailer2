package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	mailercli "github.com/tinymailer/mailer/cli"
	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/version"
)

var (
	// cli flags defination
	globalFlags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "authkey",
			Usage:  "the authentication key",
			Value:  "",
			EnvVar: "AUTH_KEY",
		},
		cli.StringFlag{
			Name:   "dburl",
			Usage:  "The connection URL of mongodb server. eg: mongodb://127.0.0.1:27017",
			Value:  "mongodb://127.0.0.1:27017/",
			EnvVar: "DB_URL",
		},
		cli.StringFlag{
			Name:   "dbname",
			Usage:  "The name of the database",
			Value:  "mailer",
			EnvVar: "DB_NAME",
		},
	}

	daemonFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "listen",
			Usage:  "The address that API service listens on",
			EnvVar: "LISTEN_ADDR",
			Value:  ":80",
		},
	}

	loadServerFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "smtp server host address. eg: smtp.126.com",
		},
		cli.StringFlag{
			Name:  "port",
			Usage: "smtp server host port. eg: 25",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "smtp server auth user email. eg: xxx@yyy.zzz",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "smtp server auth user password.",
		},
	}

	loadRecipientFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "name of recipient",
		},
		cli.StringFlag{
			Name:  "emails",
			Usage: "email list of recipient (split by ,)",
		},
	}

	loadMailFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "from-name",
			Usage: "display name of sender",
		},
		cli.StringFlag{
			Name:  "subject",
			Usage: "mail subject",
		},
		cli.StringFlag{
			Name:  "body",
			Usage: "mail content. eg: `/path/to/ad.txt` or `mail content string`",
		},
	}

	removeFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "bson id for the to be removed object",
		},
	}

	taskCreateFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "recipient",
			Usage: "bson id of recipient",
		},
		cli.StringFlag{
			Name:  "servers",
			Usage: "bson ids of server (split by ,)",
		},
		cli.StringFlag{
			Name:  "mails",
			Usage: "bson ids of mail (split by ,)",
		},
	}
	taskRunFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "bson id for the specified task",
		},
	}
	taskShowFlags   = taskRunFlags
	taskFollowFlags = taskRunFlags
	taskStopFlags   = taskRunFlags
	taskRemoveFlags = taskRunFlags
)

func main() {
	app := cli.NewApp()
	app.Name = "mailer"
	app.Usage = "simple smtp mailer"
	app.Author = ""
	app.Email = ""
	app.Version = version.GetVersion()
	if gitCommit := version.GetGitCommit(); gitCommit != "" {
		app.Version += "-" + gitCommit
	}

	app.Flags = globalFlags

	app.Before = func(c *cli.Context) error {
		var (
			debug = c.Bool("debug")
		)

		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		log.SetOutput(os.Stdout)
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "daemon",
			ShortName: "d",
			Usage:     "start the mailer daemon",
			Flags:     daemonFlags,
			Action: func(c *cli.Context) {
				var (
					listen = c.String("listen")
				)
				if err := initSetUp(c); err != nil {
					log.Fatalln(err)
				}
				if err := mailercli.Daemon(listen); err != nil {
					log.Fatalln(err)
				}
			},
		},
		cli.Command{
			Name:      "load",
			ShortName: "l",
			Usage:     "load objects",
			Before: func(c *cli.Context) error {
				return initSetUp(c)
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "server",
					Usage: "load smtp server configs",
					Flags: loadServerFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Load("server", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "recipient",
					Usage: "load recipients list",
					Flags: loadRecipientFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Load("recipient", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "mail",
					Usage: "load mail content",
					Flags: loadMailFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Load("mail", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
			},
		},
		cli.Command{
			Name:      "show",
			ShortName: "s",
			Usage:     "show objects",
			Before: func(c *cli.Context) error {
				return initSetUp(c)
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "server",
					Usage: "show smtp server configs",
					Action: func(c *cli.Context) {
						if err := mailercli.Show("server", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "recipient",
					Usage: "show recipient list",
					Action: func(c *cli.Context) {
						if err := mailercli.Show("recipient", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "mail",
					Usage: "show mail list",
					Action: func(c *cli.Context) {
						if err := mailercli.Show("mail", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
			},
		},
		cli.Command{
			Name:      "rm",
			ShortName: "d",
			Usage:     "remove objects",
			Before: func(c *cli.Context) error {
				return initSetUp(c)
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "server",
					Usage: "remove smtp server",
					Flags: removeFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Remove("server", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "recipient",
					Usage: "remove recipient",
					Flags: removeFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Remove("recipient", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "mail",
					Usage: "remove mail",
					Flags: removeFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Remove("mail", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
			},
		},
		cli.Command{
			Name:      "task",
			ShortName: "t",
			Usage:     "task manage",
			Before: func(c *cli.Context) error {
				return initSetUp(c)
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "create",
					Usage: "create a send task",
					Flags: taskCreateFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("create", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "run",
					Usage: "run a send task",
					Flags: taskRunFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("run", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "show",
					Usage: "show status of a send task",
					Flags: taskShowFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("show", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "follow",
					Usage: "follow progress of a send task",
					Flags: taskFollowFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("follow", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "stop",
					Usage: "stop a running send task",
					Flags: taskStopFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("stop", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
				cli.Command{
					Name:  "rm",
					Usage: "remove a send task",
					Flags: taskRemoveFlags,
					Action: func(c *cli.Context) {
						if err := mailercli.Task("rm", c); err != nil {
							log.Fatalln(err)
						}
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func initSetUp(c *cli.Context) error {
	var (
		dburl  = c.GlobalString("dburl")
		dbname = c.GlobalString("dbname")
	)

	// setup db connection
	return db.SetUp(dburl, dbname)
}
