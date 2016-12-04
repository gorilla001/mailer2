package db

import (
	"time"

	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	db *Database

	// DBName is exported
	DBName string
)

// Database is exported
type Database struct {
	session *mgo.Session
	dbName  string
}

// SetUp initilizing mongodb connection
func SetUp(dbURL, dbName string) (err error) {
	DBName = dbName
	if DBName == "" {
		DBName = "mailer"
	}
	if dbURL == "" {
		dbURL = "mongodb://127.0.0.1/"
	}
	db, err = NewDatabase(dbURL)
	if err != nil {
		return err
	}
	return createIndex()
}

// DB return the global initilized Database
func DB() *Database {
	return db
}

// NewDatabase creates a new database session with specified mongo URL
func NewDatabase(dbURL string) (*Database, error) {
	var (
		session *mgo.Session
		err     error
	)

	for i := 0; i <= 10; i++ {
		session, err = mgo.DialWithTimeout(dbURL, time.Second*3)
		if err == nil {
			break
		}
		log.Warnf("dial mongodb %s error: %v, retrying ...", dbURL, err)
		time.Sleep(time.Second * time.Duration(i))
	}

	if err != nil {
		return nil, err
	}
	return &Database{
		session: session,
		dbName:  DBName,
	}, nil
}

// DBName is exported
func (d *Database) DBName() string {
	return d.dbName
}

// NewSession is exported
func (d *Database) NewSession() *mgo.Session {
	return d.session.Clone()
}

// Exec is exported
func (d *Database) Exec(execHandler func(db *mgo.Database) error) error {
	var (
		ss = d.session.Clone()
		db = ss.DB(d.dbName)
	)
	defer ss.Close()
	return execHandler(db)
}

// All is exported
func (d *Database) All(collection string, query bson.M, result interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		return db.C(collection).Find(query).All(result)
	})
}

// One is exported
func (d *Database) One(collection string, query bson.M, result interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		return db.C(collection).Find(query).One(result)
	})
}

// SelectOne is exported
func (d *Database) SelectOne(collection string, query bson.M, selectQuery bson.M, result interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		return db.C(collection).Find(query).Select(selectQuery).One(result)
	})
}

// RemoveAll is exported
func (d *Database) RemoveAll(collection string, query bson.M) error {
	return d.Exec(func(db *mgo.Database) error {
		_, err := db.C(collection).RemoveAll(query)
		return err
	})
}

// Insert is exported
func (d *Database) Insert(collection string, value interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		return db.C(collection).Insert(value)
	})
}

// Update is exported
func (d *Database) Update(collection string, query bson.M, value interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		return db.C(collection).Update(query, value)
	})
}

// UpdateAll is exported
func (d *Database) UpdateAll(collection string, query bson.M, value interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		_, err := db.C(collection).UpdateAll(query, value)
		return err
	})
}

// Upsert is exported
func (d *Database) Upsert(collection string, query bson.M, value interface{}) error {
	return d.Exec(func(db *mgo.Database) error {
		_, err := db.C(collection).Upsert(query, value)
		return err
	})
}

// Ping is exported
func (d *Database) Ping() error {
	return d.session.Ping()
}

// Close close the mongo database session
func (d *Database) Close() {
	d.session.Close()
}

// BSONIDQuery is exported
func BSONIDQuery(id bson.ObjectId) bson.M {
	return bson.M{"_id": id}
}
