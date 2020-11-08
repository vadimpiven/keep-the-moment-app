// This package assists with accessing postgres and retrieving the data.
package postgres

import (
	"fmt"

	"github.com/go-pg/pg/v10"
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
	// DB is a pg.DB wrapper which makes accessing user-implemented methods easier.
	DB pg.DB
)

// New returns new DB instance.
func New(c Config) *DB {
	db := (*DB)(pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:            c.Username,
		Password:        c.Password,
		Database:        c.Database,
		ApplicationName: "go-REST",
	}))
	if err := db.CreateSchema(); err != nil {
		panic(err)
	}
	fmt.Printf("⇨ db connection established on [%s]:%d\n", c.Host, c.Port)
	return db
}

// Inject injects `db` variable in echo context.
func (db *DB) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}
