package mgodriver

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MongoDbSuite struct{ *MongoDb }

var _ = Suite(&MongoDbSuite{})

type User struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string
}

func (s *MongoDbSuite) SetUpSuite(c *C) {
	// session, err := mgo.Dial("localhost")
	// c.Check(err, Equals, nil)
	// s.MongoDb = NewMongoDb(session.DB("gogetter"))
	s.MongoDb = NewMongoDb(getTestDb())
	_, err := s.db.C("mongousers").RemoveAll(nil)
	c.Check(err, Equals, nil)
}

func (s *MongoDbSuite) TearDownSuite(c *C) {
	_, err := s.db.C("mongousers").RemoveAll(nil)
	c.Check(err, Equals, nil)
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
	err := s.Create("mongousers", user)
	c.Check(err, Equals, nil)
	count, err := s.db.C("mongousers").Find(bson.M{"name": "a user"}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 1)

	s.Remove("mongousers", user.Id)
	count, err = s.db.C("mongousers").Find(bson.M{"name": "a user"}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}
