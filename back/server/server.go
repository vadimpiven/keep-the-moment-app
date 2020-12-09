// This package contains the main part of the program which is REST api implementation.
package server

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/FTi130/keep-the-moment-app/back/lib/coords"
	"github.com/FTi130/keep-the-moment-app/back/server/api"
)

type (
	// Flags structure contains variables, which should be received as commandline flags.
	Flags struct {
		Debug *bool
	}
	// Config structure contains configurable options of this package.
	Config struct {
		Host string
		Port int
	}
	// Server structure contains configuration, commandline flags and router instance.
	Server struct {
		conf  Config
		flags Flags
		echo  *echo.Echo
	}
)

// New returns new Server instance.
func New(f Flags, c Config) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Debug = *f.Debug
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
			Skipper: func(c echo.Context) bool {
				if strings.HasPrefix(c.Request().URL.EscapedPath(), "/api/swagger") {
					return true
				}
				return false
			},
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		}),
		coords.Middleware(),
	)

	return &Server{c, f, e}
}

// Use function proxies user-implemented middleware to the internal router.
func (r *Server) Use(middleware ...echo.MiddlewareFunc) {
	r.echo.Use(middleware...)
}

// Run starts the server and loops endlessly.
func (r *Server) Run() {
	api.ApplyRoutes(r.echo)
	if *r.flags.Debug {
		fmt.Println("  starting server in DEBUG mode")
		r.echo.Logger.Debug(r.echo.Start(fmt.Sprintf("%s:%d", r.conf.Host, r.conf.Port)))
	} else {
		r.echo.Logger.Fatal(r.echo.Start(fmt.Sprintf("%s:%d", r.conf.Host, r.conf.Port)))
	}
}
