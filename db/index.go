package db

import (
	mgo "gopkg.in/mgo.v2"
)

type dbIndex struct {
	collection string
	indexes    []mgo.Index
}

var (
	indexes = []dbIndex{
		dbIndex{
			collection: CSERVER,
			indexes: []mgo.Index{
				mgo.Index{
					Key:    []string{"host", "port", "auth_user"},
					Unique: true,
				},
			},
		},
		dbIndex{
			collection: CRECIPIENT,
			indexes: []mgo.Index{
				mgo.Index{
					Key:    []string{"name"},
					Unique: true,
				},
			},
		},
	}
)

func createIndex() error {
	return DB().Exec(func(db *mgo.Database) error {
		for _, idx := range indexes {
			for _, mgoidx := range idx.indexes {
				if err := db.C(idx.collection).EnsureIndex(mgoidx); err != nil {
					return err
				}
				if err := db.C(idx.collection).EnsureIndexKey("_id"); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
