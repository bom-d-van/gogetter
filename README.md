gogetter
========

A simple factory girl port with a simple database cleaner interface, in Golang.

Gogetter is designed to be a testing tool.

# Usage

```golang
type User struct {
	Id    bson.ObjectId `bson:"_id"`
	Name  string
}

func init() {
	SetGoal("User", func() Dream {
		return User{
			Id: bson.NewObjectId(),
			Name: "name",
		}
	})
}

func TestUser(t *testing.T) {
	// Grow is Build in factory girl
	userI, err := gogetter.Grow("User")
	user := user.(User)
	user.Name == "name" // true

	// Use asterisk operator (*)
	userI, err := gogetter.Grow("*User", gogetter.Lession{
		"Name": "Custom Name",
	})
	user := userI.(*User)
	user.Name == "Custom Name" // true

	// Use Lesson to override default values
	userI, err := gogetter.Grow("User", gogetter.Lession{
		"Name": "Custom Name",
	})
	user := userI.(User)
	user.Name == "Custom Name" // true

	// Use Multiple Lessons
	usersI, err := gogetter.Grow("User", gogetter.Lession{
		"Name": "Custom Name",
	}, gogetter.Lession{
		"Name": "Another Lession",
	})
	users := usersI.([]User)
	users[0].Name == "Custom Name" // true
	users[1].Name == "Another Lession" // true

	// Realize is Create in factory girl, it will insert record(s) in a provided database
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	gogetter.SetDefaultGetterDb(session.DB("gogetter"))
	userI, err := gogetter.Realize("User")
	user := user.(User)
	user.Name == "name" // true
	// By Default, gogetter will use the lowercase, plural, and also replacing
	// spaces with underscores of the name defined with SetGoal
	users := []User{}
	session.DB("gogetter").C("users").FindId(user.Id).All(&users)
	users[0].Id == user.Id // true

	// Everything work in Grow work with Realize, except it will make changes in a database
	gogetter.Realize("User", Lession{...}, Lession{...})
	gogetter.Realize("*User", Lession{...}, Lession{...})
	...

	// Use AllInVain or Apocalypse to destory objects, i.e., remove records from databases
	gogetter.AllInVain("Users") // Will destory all "Users" just created
	gogetter.AllInVain("Users", user1, user2) // Only destory user1 and user2
	gogetter.Apocalypse("Users") // Equals to gogetter.AllInVain("Users")
	gogetter.Apocalypse("Users", "Another Goals") // Will Destroy all records of both "Users" and "Another Goals"
	gogetter.Apocalypse() // Will destroy every records

	// Of course, in most serious cases, you could use your own gogetter instead of the default one
	getter := gogetter.NewGoGetter(yourDb)
}


```