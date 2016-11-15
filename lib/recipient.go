package lib

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/db"
	"github.com/tinymailer/mailer/types"
)

// AddRecipient is exported
func AddRecipient(r types.Recipient) error {
	return db.DB().Insert(db.CRECIPIENT, r)
}

// ListRecipient is exported
func ListRecipient() ([]types.Recipient, error) {
	ret := make([]types.Recipient, 0)
	err := db.DB().All(db.CRECIPIENT, nil, &ret)
	return ret, err
}

// GetRecipient is exported
func GetRecipient(id bson.ObjectId) (types.Recipient, error) {
	var ret types.Recipient
	err := db.DB().One(db.CRECIPIENT, db.BSONIDQuery(id), &ret)
	return ret, err
}

// DelRecipient is exported
func DelRecipient(id bson.ObjectId) error {
	err := db.DB().RemoveAll(db.CRECIPIENT, db.BSONIDQuery(id))
	if err == mgo.ErrNotFound {
		return nil
	}
	return err
}
