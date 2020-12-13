// This package implements router group `api`.
package api

import (
	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/server/api/auth"
	"github.com/FTi130/keep-the-moment-app/back/server/api/hashtag"
	"github.com/FTi130/keep-the-moment-app/back/server/api/image"
	"github.com/FTi130/keep-the-moment-app/back/server/api/post"
	"github.com/FTi130/keep-the-moment-app/back/server/api/swagger"
	"github.com/FTi130/keep-the-moment-app/back/server/api/user"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		auth.ApplyRoutes(api)
		swagger.ApplyRoutes(api)
		user.ApplyRoutes(api)
		hashtag.ApplyRoutes(api)
		image.ApplyRoutes(api)
		post.ApplyRoutes(api)
	}
}
