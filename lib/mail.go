package lib

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/types"
)

// AddMail is exported
func AddMail(m types.Mail) error {
	return db.DB().Insert(db.CMAIL, m)
}

// ListMail is exported
func ListMail() ([]types.Mail, error) {
	ret := make([]types.Mail, 0)
	err := db.DB().All(db.CMAIL, nil, &ret)
	return ret, err
}

// GetMail is exported
func GetMail(id bson.ObjectId) (types.Mail, error) {
	var ret types.Mail
	err := db.DB().One(db.CMAIL, db.BSONIDQuery(id), &ret)
	return ret, err
}

// DelMail is exported
func DelMail(id bson.ObjectId) error {
	err := db.DB().RemoveAll(db.CMAIL, db.BSONIDQuery(id))
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}
