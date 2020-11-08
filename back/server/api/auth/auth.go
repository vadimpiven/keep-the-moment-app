// This package implements router group `auth`.
package auth

import (
	_ "fmt"
	"net/http"
	_ "time"

	_ "github.com/dgrijalva/jwt-go"
	_ "github.com/dimuska139/go-email-normalizer"
	_ "github.com/jakehl/goid"
	"github.com/labstack/echo/v4"
	_ "github.com/matcornic/hermes/v2"
	_ "github.com/sethvargo/go-password/password"

	"github.com/FTi130/keep-the-moment-app/back/lib/jwtclaims"
	_ "github.com/FTi130/keep-the-moment-app/back/lib/mail"
	_ "github.com/FTi130/keep-the-moment-app/back/postgres"
	_ "github.com/FTi130/keep-the-moment-app/back/redis"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group, s []byte) {
	auth := g.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/logout", logout, jwtclaims.Middleware(s))
	}
}

// login implements login procedure.
func login(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

// logout implements logout procedure.
func logout(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

// register implements registration of new user.
func register(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}
