package gogetter

import (
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type GoGetterSuite struct{}

var _ = Suite(&GoGetterSuite{})

type User struct {
	Id    bson.ObjectId `bson:"_id"`
	Name  string
	Dream *DreamS

	VisitedPlaces   []string
	ThereGreatIdeas [3]string
}

func (u User) Identity() interface{} {
	return u.Id
}

type DreamS struct {
	Title   string
	Content string
}

func init() {
	SetGoal("User", func() Dream { return makeUser() })
	SetGoal("Pointer User", func() Dream {
		user := makeUser()
		return &user
	})
}

func makeUser() User {
	return User{
		Name:            "name",
		Dream:           &DreamS{Title: "My Dream"},
		VisitedPlaces:   []string{"New York City", "San Franciso"},
		ThereGreatIdeas: [3]string{"Meet the one I love", "Build a greatest Product of all time", "Live and die peacefully"},
	}
}

func (s *GoGetterSuite) TestCommonGet(c *C) {
	user, err := Grow("User")
	c.Check(err, Equals, nil)
	c.Check(user.(User).Name, Equals, "name")
	c.Check(user.(User).Dream.Title, Equals, "My Dream")
}

func (s *GoGetterSuite) TestCommonGetWithLesson(c *C) {
	user, err := Grow("User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	})
	c.Check(err, Equals, nil)
	c.Check(user.(User).Name, Equals, "New Name")
	c.Check(user.(User).Dream.Title, Equals, "Conquer the world")
}

func (s *GoGetterSuite) TestCommonGetPointer(c *C) {
	user, err := Grow("Pointer User")
	c.Check(err, Equals, nil)
	c.Check(user.(*User).Name, Equals, "name")
	c.Check(user.(*User).Dream.Title, Equals, "My Dream")
}

func (s *GoGetterSuite) TestCommonGetPointerWithLesson(c *C) {
	user, err := Grow("Pointer User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	})
	c.Check(err, Equals, nil)
	c.Check(user.(*User).Name, Equals, "New Name")
	c.Check(user.(*User).Dream.Title, Equals, "Conquer the world")
}

func (s *GoGetterSuite) TestPointerlyGet(c *C) {
	user, err := Grow("*User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	})
	c.Check(err, Equals, nil)
	c.Check(user.(*User).Name, Equals, "New Name")
	c.Check(user.(*User).Dream.Title, Equals, "Conquer the world")
}

func (s *GoGetterSuite) TestPointerlyGetPointer(c *C) {
	user, err := Grow("*Pointer User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	})
	c.Check(err, Equals, nil)
	c.Check((**user.(**User)).Name, Equals, "New Name")
	c.Check((**user.(**User)).Dream.Title, Equals, "Conquer the world")
}

func (s *GoGetterSuite) TestGetWithLessons(c *C) {
	usersI, err := Grow("User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	}, Lesson{
		"Name": "No.2",
	})
	c.Check(err, Equals, nil)
	users, ok := usersI.([]User)
	c.Check(ok, Equals, true)
	c.Check(users[0].Name, Equals, "New Name")
	c.Check(users[0].Dream.Title, Equals, "Conquer the world")
	c.Check(users[1].Name, Equals, "No.2")
}

func (s *GoGetterSuite) TestGetWithLessonsInPointer(c *C) {
	usersI, err := Grow("*User", Lesson{
		"Name": "New Name",
		"Dream": &DreamS{
			Title: "Conquer the world",
		},
	}, Lesson{
		"Name": "No.2",
	})
	c.Check(err, Equals, nil)
	users, ok := usersI.([]*User)
	c.Check(ok, Equals, true)
	c.Check(users[0].Name, Equals, "New Name")
	c.Check(users[0].Dream.Title, Equals, "Conquer the world")
	c.Check(users[1].Name, Equals, "No.2")
}

func (s *GoGetterSuite) TestGetArraysAndSlicesInPointer(c *C) {
	usersI, err := Grow("*Pointer User", Lesson{
		"VisitedPlaces": []string{"New York City"},
	}, Lesson{
		"ThereGreatIdeas": [3]string{"1", "2", "3"},
	})
	users, ok := usersI.([]**User)
	c.Check(ok, Equals, true)
	c.Check(err, Equals, nil)
	c.Check((**users[0]).VisitedPlaces, DeepEquals, []string{"New York City"})
	c.Check((**users[0]).ThereGreatIdeas, DeepEquals, [3]string{"Meet the one I love", "Build a greatest Product of all time", "Live and die peacefully"})
	c.Check((**users[1]).VisitedPlaces, DeepEquals, []string{"New York City", "San Franciso"})
	c.Check((**users[1]).ThereGreatIdeas, DeepEquals, [3]string{"1", "2", "3"})
}

func (s *GoGetterSuite) TestCollectDreamsCorrectly(c *C) {
	gg := NewGoGetter(nil)
	_, err := gg.Grow("User")
	c.Check(err, Equals, nil)
	_, err = gg.Grow("User", nil, nil, nil)
	c.Check(err, Equals, nil)
	c.Check(gg.dreams["User"], HasLen, 4)
	// gg.Realize(name, ...)
}

func (s *GoGetterSuite) TestRealize(c *C) {
	db := getTestDb()
	SetDefaultGetterDb(NewMongoDb(db))
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
	users, ok := usersI.([]User)
	c.Check(ok, Equals, true)

	usersInDb := []User{}
	err = db.C("users").FindId(bson.M{"$in": []bson.ObjectId{users[0].Id, users[1].Id}}).All(&usersInDb)
	c.Check(err, Equals, nil)
	c.Check(usersInDb[0].Id, Equals, users[0].Id)
	c.Check(usersInDb[1].Id, Equals, users[1].Id)

	_, err = db.C("users").RemoveAll(bson.M{"_id": bson.M{"$in": []bson.ObjectId{users[0].Id, users[1].Id}}})
	c.Check(err, Equals, nil)
}

// func (s *GoGetterSuite) TestGetWithInspiration(c *C) {
// 	user, err := Grow("*Pointer User", Lesson{
// 		"Name": func() Dream {
// 			return "Name Filled by Inspiration"
// 		},
// 	})
// 	c.Check(err, Equals, nil)
// 	c.Check((**user.(**User)).Name, Equals, "Name Filled by Inspiration")
// 	c.Check((**user.(**User)).Dream.Title, Equals, "Conquer the world")
// }
