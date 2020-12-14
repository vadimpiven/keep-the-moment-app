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
	user := g.Group("/user")
	{
		user.GET("/info", getInfo, keyauth.Middleware())
		user.POST("/info", updateInfo, keyauth.Middleware())
		user.POST("/lookup", lookup)
	}
}

type (
	lookupIn struct {
		UserID string `json:"user_id"`
	}
	lookupOut200 struct {
		UserIDs []string `json:"user_ids"`
	}
)

// @Summary Get the list of hashtags similar to one that user tries to enter.
// @Accept json
// @Param user_id body lookupIn true "user_id beginning"
// @Success 200 {object} lookupOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /user/lookup [post]
func lookup(c echo.Context) error {
	in := new(lookupIn)
	err := c.Bind(in)
	if err != nil || in.UserID == "" {
		return echo.ErrBadRequest
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	if re.Find([]byte(in.UserID)) != nil {
		return echo.ErrBadRequest
	}

	userIDs, err := postgres.GetUserIDsBeginningWith(c, in.UserID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, lookupOut200{userIDs})
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
	if re.Find([]byte(user.ID)) != nil {
		return echo.ErrBadRequest
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
