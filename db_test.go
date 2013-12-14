package gogetter

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
)

type MongoDbSuite struct{ *MongoDb }

var _ = Suite(&MongoDbSuite{})

func (s *MongoDbSuite) SetUpSuite(c *C) {
	// session, err := mgo.Dial("localhost")
	// c.Check(err, Equals, nil)
	// s.MongoDb = NewMongoDb(session.DB("gogetter"))
	s.MongoDb = NewMongoDb(getTestDb())
}

func getTestDb() *mgo.Database {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return session.DB("gogetter")
}

func (s *MongoDbSuite) TestDbOperation(c *C) {
	user := User{Id: bson.NewObjectId(), Name: "a user"}
	err := s.Create("users", user)
	c.Check(err, Equals, nil)
	count, err := s.db.C("users").Find(bson.M{"name": "a user"}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 1)

	s.Remove("users", user)
	count, err = s.db.C("users").Find(bson.M{"name": "a user"}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}
