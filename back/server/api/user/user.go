// This package implements router group `user`.
package user

import (
	"net/http"
	"regexp"

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

// @Summary Return information about user, or 404 if user not registered.
// @Security Bearer
// @Produce json
// @Success 200 {object} postgres.User
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /user/info [get]
func getInfo(c echo.Context) error {
	email, err := keyauth.GetEmail(c)
	if err != nil {
		return echo.ErrInternalServerError
	}

	user, err := postgres.GetUser(c, email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Updates information about user.
// @Security Bearer
// @Accept json
// @Produce json
// @Param user_info body postgres.User true "all information about user"
// @Success 200 {object} postgres.User
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /user/info [post]
func updateInfo(c echo.Context) error {
	user := new(postgres.User)
	err := c.Bind(user)
	if err != nil {
		return echo.ErrBadRequest
	} else if user.Hashtags == nil {
		user.Hashtags = []string{}
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	for _, hashtag := range user.Hashtags {
		if re.Find([]byte(hashtag)) != nil {
			return echo.ErrBadRequest
		}
	}

	valid, err := postgres.CheckUserValid(c, user)
	if err != nil {
		return echo.ErrInternalServerError
	} else if valid == false {
		return echo.ErrBadRequest
	}

	err = postgres.UpdateUser(c, user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, user)
}
