// This package implements key-auth validation middleware.
package keyauth

import (
	"errors"
	"net/http"
	"time"

	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
	"github.com/labstack/echo/v4"
)

const authScheme = "Bearer"

// Middleware returns preconfigured jwt middleware.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := getToken(c)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "auth token not provided",
					Internal: err,
				}
			}
			valid, err := redis.CheckTokenExistsAndProlong(c, token, 72*time.Hour)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "error while checking auth token validity",
					Internal: err,
				}
			} else if valid {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func getToken(c echo.Context) (string, error) {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	if auth == "" {
		return "", errors.New("missing key in request header")
	}
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", errors.New("invalid key in the request header")
}

func GetEmail(c echo.Context) (string, error) {
	token, err := getToken(c)
	if err != nil {
		return "", err
	}
	return redis.GetValueByToken(c, token)
}

func ExpireToken(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	return redis.DeleteToken(c, token)
}
