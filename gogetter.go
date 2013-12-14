package gogetter

import (
	"bitbucket.org/pkg/inflect"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Dream interface {
	// Identity() interface{} // Id
	// EqualTo(anotherDream Dream) bool
}
type Goal func() Dream
type Lesson map[string]Dream

// TODO: Support Special Type
// type Inspiration func() Dream

type GoGetter struct {
	dreams map[string][]Dream
	db     Database
}

func NewGoGetter(db Database) *GoGetter {
	return &GoGetter{
		db:     db,
		dreams: map[string][]Dream{},
	}
}

var ErrGetterNotExist = errors.New("Getter Not Exist")
var defaultGetter = NewGoGetter(nil)
var goalMap = map[string]Goal{}
var tableNameMap = map[string]string{}

// var dreamType = reflect.TypeOf()

// Setting table name is optional, if table name is not specifically setted, gogetter
// will use the pluralization and lower case form of the name as table name
func SetTableName(name, table string) {
	tableNameMap[name] = table
	// tableNameMap["*"+name] = table
}

func GetTableName(name string) (table string, err error) {
	// if len(name) > 0 && name[0] == '*' {
	// 	name = name[1:]
	// }

	var ok bool
	table, ok = tableNameMap[name]
	if ok {
		return
	}

	if GetGoal(name) == nil {
		return "", ErrGetterNotExist
	}

	table = inflect.Pluralize(strings.ToLower(name))
	tableNameMap[name] = table

	return
}

var mux = sync.Mutex{}

// Note: Leading asterisk (*) in name is saved for gogetter.
func SetGoal(name string, goal Goal) {
	mux.Lock()
	defer mux.Unlock()

	goalMap[name] = goal
	// goalMap["*"+name] = goal
}

func GetGoal(name string) Goal {
	// if len(name) > 0 && name[0] == '*' {
	// 	name = name[1:]
	// }

	goal, ok := goalMap[name]
	if !ok {
		return nil
	}

	return goal
}

func SetDefaultGetterDb(db Database) {
	defaultGetter.db = db
}

// func Raise(children, parent string, lessons ...Lesson) {
// }

func Grow(name string, lessons ...Lesson) (Dream, error) {
	return defaultGetter.Grow(name, lessons...)
}

// Must calling SetDefaultGetterDb before this
func Realize(name string, lessons ...Lesson) (dreams Dream, err error) {
	return defaultGetter.Realize(name, lessons...)
}

// Must calling SetDefaultGetterDb before this
func AllInVain(name string, dreams ...Dream) (err error) {
	return defaultGetter.AllInVain(name, dreams...)
}

// Could use a leading asterisk (*) in name to get pointer value
// TODO:
// 	1. Support anonymous type, e,g, custom struct, map, etc
// 	2. Tags specification in custom structs, provided that gogetter will support struct
func (gg *GoGetter) Grow(name string, lessons ...Lesson) (dreams Dream, err error) {
	return gg.makeDreams(name, false, lessons...)
}

func (gg *GoGetter) makeDreams(name string, saveInDb bool, lessons ...Lesson) (dreams Dream, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
		}
	}()

	// if _, ok := goal; !ok {
	// 	return nil, ErrGetterNotExist
	// }
	inPointer := len(name) > 1 && name[0] == '*'
	if inPointer {
		name = name[1:]
	}
	goal := GetGoal(name)
	if goal == nil {
		return nil, ErrGetterNotExist
	}

	// Start Produce Dreams
	firstD := reflect.ValueOf(goal())
	dType := firstD.Type()
	if inPointer {
		dType = reflect.PtrTo(dType)
	}
	goalChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, dType), 1)
	if len(lessons) == 0 {
		lessons = append(lessons, nil)
	}
	for i, _ := range lessons {
		go func(index int) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%+v", r)
					goalChan.Send(reflect.Zero(dType))
				}
			}()

			var forebear reflect.Value
			if index == 0 {
				forebear = firstD
			} else {
				forebear = reflect.ValueOf(goal())
			}

			// TODO: refactor
			var dst, src reflect.Value
			theone := reflect.New(dType)
			if inPointer {
				if forebear.Kind() == reflect.Ptr {
					src = forebear.Elem()
					v := reflect.New(src.Type())
					v.Elem().Set(reflect.New(src.Type()).Elem())
					vv := reflect.New(v.Type())
					vv.Elem().Set(v)
					theone.Elem().Set(vv)
					dst = theone.Elem().Elem().Elem()
				} else {
					src = forebear
					v := reflect.New(src.Type())
					v.Elem().Set(reflect.New(src.Type()).Elem())
					theone.Elem().Set(v)
					dst = theone.Elem().Elem()
				}
			} else {
				if dType.Kind() == reflect.Ptr {
					src = forebear.Elem()
					v := reflect.New(src.Type())
					v.Elem().Set(reflect.New(src.Type()).Elem())
					theone.Elem().Set(v)
					dst = theone.Elem().Elem()
				} else {
					src = forebear
					dst = theone.Elem()
				}
			}

			for j := 0; j < src.NumField(); j++ {
				fIndex := []int{j}
				v := src.FieldByIndex(fIndex)
				dst.FieldByIndex(fIndex).Set(v)
			}

			lesson := lessons[index]
			for k, v := range lesson {
				dst.FieldByName(k).Set(reflect.ValueOf(v))
			}

			goalChan.Send(theone.Elem())
		}(i)
	}

	// Receive Dreams
	goals := reflect.MakeSlice(reflect.SliceOf(dType), 0, 0)
	for i := 0; i < len(lessons); i++ {
		goal, _ := goalChan.Recv()
		goals = reflect.Append(goals, goal)
		gg.dreams[name] = append(gg.dreams[name], goal.Interface())
	}

	table := ""
	if saveInDb {
		table, err = GetTableName(name)
		if err != nil {
			return
		}
	}

	if saveInDb {
		records := []interface{}{}
		for i := 0; i < goals.Len(); i++ {
			records = append(records, goals.Index(i).Interface())
		}
		err = gg.db.Create(table, records...)
	}

	// Return userful/handy results
	if goals.Len() == 0 {
		dreams = reflect.Zero(dType).Interface()
	} else if goals.Len() == 1 {
		dreams = goals.Index(0).Interface()
	} else {
		dreams = goals.Interface()
	}

	return
}

// Grow and Create a Record in Database
func (gg *GoGetter) Realize(name string, lessons ...Lesson) (dreams Dream, err error) {
	return gg.makeDreams(name, true, lessons...)
}

var allInVainMutex = sync.Mutex{}

// Remove from database (Do not use leading * in name with this function)
// TODO: enable field tag configuration
// TODO: !!!IMPORTANT!!! support slice type arguments
func (gg *GoGetter) AllInVain(name string, dreams ...Dream) (err error) {
	allInVainMutex.Lock()
	defer allInVainMutex.Unlock()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = fmt.Errorf("%+v", r)
	// 	}
	// }()

	if len(dreams) == 1 && reflect.TypeOf(dreams[0]).Kind() == reflect.Slice {
		packedDreams := reflect.ValueOf(dreams[0])
		dreams = []Dream{}
		for i := 0; i < packedDreams.Len(); i++ {
			dreams = append(dreams, packedDreams.Index(i).Interface())
		}
	}

	table, err := GetTableName(name)
	if err != nil {
		return
	}

	if len(dreams) == 0 {
		dreams = gg.dreams[name]
		if len(dreams) == 0 {
			return
		}
	}

	records := []Record{}
	for i, _ := range dreams {
		// TODO: use func to acess Record
		records = append(records, dreams[i].(Record))
	}

	restDreams := []Dream{}
	for _, dream := range gg.dreams[name] {
		// TODO: use func to acess Record
		dreamRecord := dream.(Record)
		for _, dyingDream := range records {
			if reflect.DeepEqual(dreamRecord.Identity(), dyingDream.Identity()) {
				goto hell
			}
		}

		restDreams = append(restDreams, dream)

	hell:
	}
	gg.dreams[name] = restDreams

	return gg.db.Remove(table, records...)
}

func (gg *GoGetter) Apocalypse(names ...string) (err error) {
	// if len(names) == 0 {
	// 	for k, _ := range gg.dreams {
	// 		// names =
	// 	}
	// }

	return
}

// func (gg *GoGetter) Raise(children, parent string, lessons ...Lesson) {
// }
