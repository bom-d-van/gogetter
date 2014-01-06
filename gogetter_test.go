package gogetter

import (
	"github.com/bom-d-van/gogetter/mgodriver"
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

// func (u User) Identity() interface{} {
// 	return u.Id
// }

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

	AscendGoal("Super User", "User", func() Lesson {
		return Lesson{
			"Name": "Super User",
		}
	})
}

func makeUser() User {
	return User{
		Id:              bson.NewObjectId(),
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
	c.Check(users[1].Dream.Title, Equals, "My Dream")
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
	SetDefaultGetterDb(mgodriver.NewMongoDb(db))
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

func (s *GoGetterSuite) TestGetTableName(c *C) {
	SetGoal("name with space", func() Dream { return nil })
	SetGoal("Capital", func() Dream { return nil })
	SetGoal("Capital  and Space", func() Dream { return nil })

	var table string
	var err error
	table, err = GetTableName("name with space")
	c.Check(err, Equals, nil)
	c.Check(table, Equals, "name_with_spaces")
	table, err = GetTableName("Capital")
	c.Check(err, Equals, nil)
	c.Check(table, Equals, "capitals")
	table, err = GetTableName("Capital  and Space")
	c.Check(err, Equals, nil)
	c.Check(table, Equals, "capital__and_spaces")
}

func (s *GoGetterSuite) TestRetrieveDreamId(c *C) {
	var id interface{}
	puser := User{}
	id = defaultGetter.retrieveDreamId(puser, "Id")
	c.Check(id, Not(Equals), nil)
	id = defaultGetter.retrieveDreamId(&puser, "Id")
	c.Check(id, Not(Equals), nil)
	ppuser := &User{}
	id = defaultGetter.retrieveDreamId(&ppuser, "Id")
	c.Check(id, Not(Equals), nil)
}

func (s *GoGetterSuite) TestGetDreamIdField(c *C) {
	cidCalledCount := 0
	SetGoal("CustomId", func() Dream {
		cidCalledCount += 1
		return struct {
			CustomId string `gogetter:"id"`
		}{}
	})

	c.Check(getDreamIdField("CustomId"), Equals, "CustomId")

	// Should cached DreamId
	getDreamIdField("CustomId")
	getDreamIdField("CustomId")
	c.Check(cidCalledCount, Equals, 1)

	SetGoal("WithOutId", func() Dream {
		return struct {
		}{}
	})

	c.Check(getDreamIdField("WithOutId"), Equals, "")
}

func (s *GoGetterSuite) TestGetTableNameOfAscendGoals(c *C) {
	table, err := GetTableName("Super User")
	c.Check(err, Equals, nil)
	c.Check(table, Equals, "users")
}

func (s *GoGetterSuite) TestNestedAscendGoals(c *C) {
	AscendGoal("Super Super User", "Super User", func() Lesson {
		return Lesson{
			"Name": "Super Super User",
			"Dream": &DreamS{
				Title: "Super Super Dream",
			},
		}
	})

	userI, err := Grow("*Super Super User")
	c.Check(err, Equals, nil)
	user := userI.(*User)
	c.Check(user.Name, Equals, "Super Super User")
	c.Check(user.Dream.Title, Equals, "Super Super Dream")
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
