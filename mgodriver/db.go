package mgodriver

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDb struct {
	db *mgo.Database
}

func NewMongoDb(db *mgo.Database) (mdb *MongoDb) {
	return &MongoDb{db: db}
}

func (m *MongoDb) Create(col string, docs ...interface{}) (err error) {
	if len(docs) == 0 {
		return
	}

	return m.db.C(col).Insert(docs...)
}

// Note: If the id field is "Id", MongoDb will convert it into "_id", which is the real id of documents in mongodb.
// If this conversion doesn't fit in your cases, feel free to create you own mongo db driver.
func (m *MongoDb) Remove(col string, idField string, ids ...interface{}) (err error) {
	if idField == "Id" {
		idField = "_id"
	}
	_, err = m.db.C(col).RemoveAll(bson.M{idField: bson.M{"$in": ids}})
	return
}
