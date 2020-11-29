package postgres

import (
	"github.com/labstack/echo/v4"
)

// Query methods documented here: https://pg.uptrace.dev/queries/

// CheckUserExists checks whether user with given email already registered or not.
func CheckUserExists(c echo.Context, email string) (res bool, err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	return db.ModelContext(ctx, (*User)(nil)).
		Where("email = ?", email).
		Exists()
}

// GetUserInfo retrieves information about user from the postgres.
func GetUserInfo(c echo.Context, email string) (res User, err error) {
	db := c.Get(contextKey).(*Postgres)
	ctx := c.Request().Context()

	err = db.ModelContext(ctx, &res).
		Where("email = ?", email).
		Select()
	return
}
