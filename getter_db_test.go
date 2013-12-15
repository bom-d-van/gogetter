package gogetter

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
)

type GetterDbSuite struct {
	users   []User
	pusers  []*User
	ppusers []**User
	db      *mgo.Database
}

var _ = Suite(&GetterDbSuite{})

func init() {
	SetGoal("PUser", func() Dream { return makePUser() })
}

func makePUser() *User {
	return &User{
		Name:            "name",
		Dream:           &DreamS{Title: "My Dream"},
		VisitedPlaces:   []string{"New York City", "San Franciso"},
		ThereGreatIdeas: [3]string{"Meet the one I love", "Build a greatest Product of all time", "Live and die peacefully"},
	}
}

func (s *GetterDbSuite) SetUpSuite(c *C) {
	s.db = getTestDb()
	SetDefaultGetterDb(NewMongoDb(s.db))
}

func (s *GetterDbSuite) SetUpTest(c *C) {
	defaultGetter.dreams["User"] = nil
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
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 2)

	defaultGetter.dreams["PUser"] = nil
	pointerUserI, err := Realize("PUser", Lesson{"Id": bson.NewObjectId()}, Lesson{"Id": bson.NewObjectId()})
	c.Check(err, Equals, nil)
	s.pusers = pointerUserI.([]*User)
	count, err = s.db.C("pusers").FindId(bson.M{"$in": []bson.ObjectId{s.pusers[0].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 1)

	defaultGetter.dreams["Pointer User"] = nil
	ppuserI, err := Realize("*Pointer User", Lesson{"Id": bson.NewObjectId()}, Lesson{"Id": bson.NewObjectId()})
	c.Check(err, Equals, nil)
	s.ppusers = ppuserI.([]**User)
	count, err = s.db.C("pointer_users").FindId(bson.M{"$in": []bson.ObjectId{(*s.ppusers[0]).Id, (*s.ppusers[1]).Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 2)
}

func (s *GetterDbSuite) TearDownTest(c *C) {
	_, err := s.db.C("users").RemoveAll(nil)
	c.Check(err, Equals, nil)
	_, err = s.db.C("pusers").RemoveAll(nil)
	c.Check(err, Equals, nil)
	// _, err = s.db.C("pointer_users").RemoveAll(nil)
	// c.Check(err, Equals, nil)
}

func (s *GetterDbSuite) TestCommonAllInVain(c *C) {
	c.Check(defaultGetter.dreams["User"], HasLen, 2)
	err := AllInVain("User", s.users[0])
	c.Check(err, Equals, nil)
	c.Check(defaultGetter.dreams["User"], HasLen, 1)
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 1)
}

func (s *GetterDbSuite) TestAllInVainReceivingSlice(c *C) {
	c.Check(defaultGetter.dreams["User"], HasLen, 2)
	err := AllInVain("User", s.users)
	c.Check(err, Equals, nil)
	c.Check(defaultGetter.dreams["User"], HasLen, 0)
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

func (s *GetterDbSuite) TestAllInVainWithZeroArguments(c *C) {
	c.Check(defaultGetter.dreams["User"], HasLen, 2)
	err := AllInVain("User")
	c.Check(err, Equals, nil)
	c.Check(defaultGetter.dreams["User"], HasLen, 0)
	count, err := s.db.C("users").FindId(bson.M{"$in": []bson.ObjectId{s.users[0].Id, s.users[1].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

func (s *GetterDbSuite) TestPointerAllInVain(c *C) {
	c.Check(defaultGetter.dreams["PUser"], HasLen, 2)
	err := AllInVain("PUser", s.pusers[0])
	c.Check(err, Equals, nil)
	c.Check(defaultGetter.dreams["PUser"], HasLen, 1)
	count, err := s.db.C("pusers").FindId(bson.M{"$in": []bson.ObjectId{s.pusers[0].Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

func (s *GetterDbSuite) TestDoublePointerAllInVain(c *C) {
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 2)
	err := AllInVain("Pointer User", s.ppusers[0])
	c.Check(err, Equals, nil)
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 1)
	count, err := s.db.C("pointer_users").FindId(bson.M{"$in": []bson.ObjectId{(*s.ppusers[0]).Id}}).Count()
	c.Check(err, Equals, nil)
	c.Check(count, Equals, 0)
}

func (s *GetterDbSuite) TestApocalypseWithNames(c *C) {
	c.Check(defaultGetter.dreams["User"], HasLen, 2)
	c.Check(defaultGetter.dreams["PUser"], HasLen, 2)
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 2)

	err := Apocalypse("User", "PUser")
	c.Check(err, Equals, nil)

	c.Check(defaultGetter.dreams["User"], HasLen, 0)
	c.Check(defaultGetter.dreams["PUser"], HasLen, 0)
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 2)
}

func (s *GetterDbSuite) TestApocalypseWithoutNames(c *C) {
	c.Check(defaultGetter.dreams["User"], HasLen, 2)
	c.Check(defaultGetter.dreams["PUser"], HasLen, 2)
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 2)

	err := Apocalypse()
	c.Check(err, Equals, nil)

	c.Check(defaultGetter.dreams["User"], HasLen, 0)
	c.Check(defaultGetter.dreams["PUser"], HasLen, 0)
	c.Check(defaultGetter.dreams["Pointer User"], HasLen, 0)
}

// TODO:
// 	3. Test Realize With Pointer
