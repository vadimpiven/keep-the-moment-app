// This package implements router group `swagger`.
package swagger

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/FTi130/keep-the-moment-app/back/docs"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	g.GET("/swagger/*", echoSwagger.WrapHandler)
}
