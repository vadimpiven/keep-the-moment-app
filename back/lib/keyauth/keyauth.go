// This package implements key-auth validation middleware.
package keyauth

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
)

const authScheme = "Bearer"

// Middleware returns preconfigured jwt middleware.
func Middleware() echo.MiddlewareFunc {
	conf := middleware.KeyAuthConfig{
		AuthScheme: authScheme,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		Validator: func(key string, c echo.Context) (bool, error) {
			return redis.CheckTokenExistsAndProlong(c, key, 72*time.Hour)
		},
	}
	return middleware.KeyAuthWithConfig(conf)
}

func getToken(c echo.Context) string {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	return auth[len(authScheme)+1:]
}

func GetEmail(c echo.Context) (string, error) {
	return redis.GetValueByToken(c, getToken(c))
}

func ExpireToken(c echo.Context) error {
	return redis.DeleteToken(c, getToken(c))
}
