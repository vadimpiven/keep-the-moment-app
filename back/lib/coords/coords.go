// This package implements user coordinates saving middleware.
package coords

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"

	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/lib/keyauth"
)

// Middleware returns preconfigured jwt middleware.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skipper
			if strings.HasPrefix(c.Request().URL.EscapedPath(), "/api/auth/login") {
				return next(c)
			}

			// Middleware
			email, err := keyauth.GetEmail(c)
			if err != nil {
				return next(c)
			}

			latText := c.Request().Header.Get("Latitude")
			if latText == "" {
				return &echo.HTTPError{
					Code:    http.StatusBadRequest,
					Message: "latitude not provided",
				}
			}
			lat, err := strconv.ParseFloat(latText, 64)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "bad latitude value",
					Internal: err,
				}
			}

			lonText := c.Request().Header.Get("Longitude")
			if lonText == "" {
				return &echo.HTTPError{
					Code:    http.StatusBadRequest,
					Message: "longitude not provided",
				}
			}
			lon, err := strconv.ParseFloat(lonText, 64)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "bad longitude value",
					Internal: err,
				}
			}

			err = postgres.UpdateUserLocation(c, email, lat, lon)
			if err != nil {
				return echo.ErrInternalServerError
			}
			return next(c)
		}
	}
}
