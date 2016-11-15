package lib

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/types"
)

// AddServer is exported
func AddServer(s types.SMTPServer) error {
	return db.DB().Insert(db.CSERVER, s)
}

// ListServer is exported
func ListServer() ([]types.SMTPServer, error) {
	ret := make([]types.SMTPServer, 0)
	err := db.DB().All(db.CSERVER, nil, &ret)
	return ret, err
}

// GetServer is exported
func GetServer(id bson.ObjectId) (types.SMTPServer, error) {
	var ret types.SMTPServer
	err := db.DB().One(db.CSERVER, db.BSONIDQuery(id), &ret)
	return ret, err
}

// DelServer is exported
func DelServer(id bson.ObjectId) error {
	err := db.DB().RemoveAll(db.CSERVER, db.BSONIDQuery(id))
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}
