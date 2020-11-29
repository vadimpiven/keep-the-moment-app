// This package assists with accessing postgres and retrieving the data.
package postgres

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"
)

type (
	// Config structure contains configurable options of this package.
	Config struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
	// Postgres is a pg.DB wrapper.
	Postgres pg.DB
)

// New returns new DB instance.
func New(c Config) *Postgres {
	db := (*Postgres)(pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:            c.Username,
		Password:        c.Password,
		Database:        c.Database,
		ApplicationName: "back",
	}))

	createTestSchema := func(db *Postgres) error {
		type ConnectionTest struct{ Dummy bool }
		return db.Model((*ConnectionTest)(nil)).
			CreateTable(&orm.CreateTableOptions{Temp: true})
	}
	if err := createTestSchema(db); err != nil {
		panic(err)
	}

	fmt.Printf("⇨ db connection established on [%s]:%d\n", c.Host, c.Port)
	return db
}

const contextKey = "__postgres__"

// Inject injects DB in echo context.
func (db *Postgres) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextKey, db)
			return next(c)
		}
	}
}
