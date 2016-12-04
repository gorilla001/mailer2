package types

import (
	"gopkg.in/mgo.v2/bson"

	"fmt"
)

// Recipient is exported
type Recipient struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	Name   string        `bson:"name" json:"name"`
	Emails []string      `bson:"emails" json:"emails"`
}

// Validate is exported
func (r *Recipient) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name required")
	}
	if len(r.Emails) == 0 {
		return fmt.Errorf("emails required")
	}
	for _, e := range r.Emails {
		if len(e) < 4 {
			return fmt.Errorf("email [%s] invalid", e)
		}
	}
	return nil
}
