package hooddriver

import (
	"database/sql"
	"github.com/eaigner/hood"
	_ "github.com/go-sql-driver/mysql"
)

type Hood struct {
	hood *hood.Hood
}

func NewHood(hood *hood.Hood) *Hood {
	return &Hood{hood: hood}
}

func (m *Hood) Create(table string, records ...interface{}) (err error) {
	hd := m.hood.Begin()
	err = hd.CreateTableIfNotExists(&User{})
	if err != nil {
		return
	}
	hd.Commit()

	hd = m.hood.Begin()
	_, err = hd.SaveAll(&records)
	hd.Commit()

	return
}

func (m *Hood) Remove(table string, idField string, ids ...interface{}) (err error) {
	hd := m.hood.Begin()
	err = hd.Where(idField, "IN", ids).DeleteFrom(table)
	hd.Commit()
	return
}
