package gogetter

import (
	// "github.com/eaigner/hood"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Database interface {
	Create(table string, data ...interface{}) (err error)
	Remove(table string, ids ...interface{}) (err error)
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

func (m *MongoDb) Remove(col string, ids ...interface{}) (err error) {
	_, err = m.db.C(col).RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	return
}

// type Hood struct {
// 	hood *hood.Hood
// }

// func NewHood(hood *hood.Hood) *Hood {
// 	return &Hood{
// 		primaryKey: primaryKey,
// 		hood:       hood,
// 	}
// }

// func (m *Hood) Create(table string, records ...interface{}) (err error) {
// 	if len(records) == 0 {
// 		return
// 	}

// 	m.hood.CreateTableIfNotExists(table)

// 	return
// }

// func (m *Hood) Remove(table string, docs ...interface{}) (err error) {
// 	_, err = m.db.C(table).RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
// 	return
// }
