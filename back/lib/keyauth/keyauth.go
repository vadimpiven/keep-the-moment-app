package keyauth

import (
	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const authScheme = "Bearer"

// Middleware returns preconfigured jwt middleware.
func Middleware() echo.MiddlewareFunc {
	conf := middleware.KeyAuthConfig{
		AuthScheme: authScheme,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		Validator: func(key string, c echo.Context) (bool, error) {
			return redis.CheckTokenExists(c, key)
		},
	}
	return middleware.KeyAuthWithConfig(conf)
}

func getToken(c echo.Context) string {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	return auth[len(authScheme)+1:]
}

// GetEmail return user email by auth token.
func GetEmail(c echo.Context) (string, error) {
	return redis.GetValue(c, getToken(c))
}

// ExpireToken performs logout process.
func ExpireToken(c echo.Context) error {
	return redis.DeleteToken(c, getToken(c))
}
