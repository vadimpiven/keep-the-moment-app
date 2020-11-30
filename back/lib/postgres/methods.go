package postgres

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

// Query methods documented here: https://pg.uptrace.dev/queries/

func RegisterIfNewUser(c echo.Context, email string) (err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		exists, err := tx.ModelContext(ctx, (*User)(nil)).
			Where("email = ?", email).
			Exists()
		if err != nil {
			return err
		}
		if exists == false {
			_, err = tx.ModelContext(ctx, &User{Email: email}).
				Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetUserInfo(c echo.Context, email string) (res User, err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	err = db.ModelContext(ctx, &res).
		Where("email = ?", email).
		Select()
	return
}

func CheckUserValid(c echo.Context, val *User) (res bool, err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		res, err = tx.ModelContext(ctx, (*User)(nil)).
			Where("email = ? AND registered = ?", val.Email, val.Registered).
			Exists()
		if err != nil {
			return err
		} else if res == false {
			return nil
		}

		res, err = tx.ModelContext(ctx, (*Image)(nil)).
			Where("path = ?", val.Image).
			Exists()
		return err
	})
	return
}

func UpdateUserInfo(c echo.Context, val *User) (err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	val.Updated = time.Now()
	_, err = db.ModelContext(ctx, val).
		WherePK().
		Returning("*").
		Update()
	return
}
