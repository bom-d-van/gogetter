package gogetter

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
)

type GetterDbSuite struct {
	users []User
	db    *mgo.Database
}

var _ = Suite(&GetterDbSuite{})

func (s *GetterDbSuite) SetUpTest(c *C) {
	s.db = getTestDb()
	SetDefaultGetterDb(NewMongoDb(s.db))
	usersI, err := Realize("User", Lesson{
		"Id":   bson.NewObjectId(),
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	}, Lesson{
		"Id":   bson.NewObjectId(),
		"Name": "No.2",
	})
	c.Check(err, Equals, nil)
	s.users = usersI.([]User)
}

func (s *GetterDbSuite) TestCommonAllInVain(c *C) {
	err := AllInVain("User", s.users[0], s.users[1])
	c.Check(err, Equals, nil)
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

func (s *GetterDbSuite) TestAllInVainReceivingSlice(c *C) {
	err := AllInVain("User", s.users)
	c.Check(err, Equals, nil)
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

// TODO:
//	1. Test gg.dreams[name]
// 	2. Test Pointer AllInVain
// 	3. Test Realize With Pointer
