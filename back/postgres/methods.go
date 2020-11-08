package postgres

import (
	"github.com/go-pg/pg/v10/orm"
)

func (db *DB) CreateSchema() (err error) {
	models := []interface{}{
		(*User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}