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

func (m *MongoDb) Remove(col string, ids ...interface{}) (err error) {
	_, err = m.db.C(col).RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	return
}
