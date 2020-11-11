// This package implements router group `swagger`.
package swagger

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group, s []byte) {
	g.GET("/swagger/*", echoSwagger.WrapHandler)
}
