// This package implements router group `hashtag`.
package hashtag

import (
	"net/http"
	"regexp"

	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"

	"github.com/labstack/echo/v4"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	auth := g.Group("/hashtag")
	{
		auth.POST("/lookup", lookup)
		auth.GET("/top", top)
	}
}

type (
	lookupIn struct {
		Hashtag string `json:"hashtag"`
	}
	lookupOut200 struct {
		Hashtags []string `json:"hashtags"`
	}
)

// @Summary Get the list of hashtags similar to one that user tries to enter.
// @Accept json
// @Param hashtag body lookupIn true "hashtag name beginning"
// @Success 200 {object} lookupOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /hashtag/lookup [post]
func lookup(c echo.Context) error {
	in := new(lookupIn)
	err := c.Bind(in)
	if err != nil || in.Hashtag == "" {
		return echo.ErrBadRequest
	}

	re := regexp.MustCompile("[^a-z0-9_]+")
	if re.Find([]byte(in.Hashtag)) != nil {
		return echo.ErrBadRequest
	}

	hashtags, err := postgres.GetHashtagsBeginningWith(c, in.Hashtag)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, lookupOut200{hashtags})
}

type (
	topOut200 struct {
		Hashtags []string `json:"hashtags"`
	}
)

// @Summary Returns the global top 10 of hashtags.
// @Produce json
// @Success 200 {object} topOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /hashtag/top [get]
func top(c echo.Context) error {
	hashtags, err := postgres.GetHashtagsBeginningWith(c, "")
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, lookupOut200{hashtags})
}
