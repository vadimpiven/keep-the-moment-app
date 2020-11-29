// This package implements router group `user`.
package user

import (
	"net/http"

	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/lib/keyauth"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	auth := g.Group("/user")
	{
		auth.GET("/info", getInfo, keyauth.Middleware())
		auth.POST("/info", updateInfo, keyauth.Middleware())

	}
}

type (
	getInfoOut200 postgres.User
)

// Return information about user, or 404 if user not registered.
// @Summary Return information about user, or 404 if user not registered.
// @Produce json
// @Success 200 {object} getInfoOut200
// @Failure 400,401,404,500 {object} httputil.HTTPError
// @Router /user/info [get]
func getInfo(c echo.Context) error {
	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if exists, err := postgres.CheckUserExists(c, email); err != nil {
		return echo.ErrInternalServerError
	} else if exists == false {
		return c.NoContent(http.StatusNotFound)
	}

	user, err := postgres.GetUserInfo(c, email)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, user)
}

type (
	updateInfoIn     postgres.User
	updateInfoOut200 postgres.User
)

// Updates information about user.
// @Summary Updates information about user.
// @Accept json,mpfd
// @Param user_info body updateInfoIn true "all information about user"
// @Success 200,201 {object} getInfoOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /user/info [post]
func updateInfo(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusNotImplemented)
}
