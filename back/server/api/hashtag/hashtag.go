// This package implements router group `hashtag`.
package hashtag

import (
	"net/http"

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
		Hashtag string `json:"hashtag" form:"hashtag"`
	}
	lookupOut200 struct {
		Hashtags []string `json:"hashtags"`
	}
)

// Get the list of hashtags similar to one that user tries to enter.
// @Summary Get the list of hashtags similar to one that user tries to enter.
// @Accept json,mpfd
// @Param hashtag body lookupIn true "hashtag name beginning"
// @Success 200 {object} lookupOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /hashtag/lookup [post]
func lookup(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusNotImplemented)
}

type (
	topOut200 struct {
		Hashtags []string `json:"hashtags"`
	}
)

// Returns the global top 10 of hashtags.
// @Summary Returns the global top 10 of hashtags.
// @Produce json
// @Success 200 {object} topOut200
// @Failure 400,500 {object} httputil.HTTPError
// @Router /hashtag/top [get]
func top(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusNotImplemented)
}
