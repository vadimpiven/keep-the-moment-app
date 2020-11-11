// This package implements router group `api`.
package api

import (
	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/server/api/auth"
	"github.com/FTi130/keep-the-moment-app/back/server/api/swagger"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(e *echo.Echo, s []byte) {
	api := e.Group("/api")
	{
		auth.ApplyRoutes(api, s)
		swagger.ApplyRoutes(api, s)
	}
}
