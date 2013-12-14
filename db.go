package gogetter

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// type Indentity interface {
// 	EqualTo(Indentity) bool
// }

type Record interface {
	Identity() interface{}
}

type Database interface {
	Create(table string, data ...interface{}) (err error)
	Remove(table string, records ...Record) (err error)
}

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

func (m *MongoDb) Remove(col string, docs ...Record) (err error) {
	ids := []interface{}{}
	for _, doc := range docs {
		ids = append(ids, doc.Identity())
	}
	_, err = m.db.C(col).RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	return
}
