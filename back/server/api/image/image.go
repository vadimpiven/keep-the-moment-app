// This package implements router group `image`.
package image

import (
	"bytes"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/lib/closable"
	"github.com/FTi130/keep-the-moment-app/back/lib/keyauth"
	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	auth := g.Group("/image")
	{
		auth.POST("/upload", upload, keyauth.Middleware())
		auth.POST("/post", post, keyauth.Middleware())
	}
}

type (
	uploadOut200 struct {
		Image string `json:"image"`
	}
)

// @Summary Updates information about user.
// @Security Bearer
// @Accept mpfd
// @Produce json
// @Param image formData file true "image file"
// @Success 200 {object} uploadOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /image/upload [post]
func upload(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil || file.Filename == "" {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "image is not attached",
			Internal: err,
		}
	}

	src, err := file.Open()
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer closable.SafeClose(src)

	scn := io.Reader(src)
	img, err := imaging.Decode(scn, imaging.AutoOrientation(true))
	if err != nil {
		return echo.ErrInternalServerError
	}

	img = imaging.Fill(img, 100, 100, imaging.Center, imaging.CatmullRom)

	buf := &bytes.Buffer{}
	if err = imaging.Encode(buf, img, imaging.PNG); err != nil {
		return echo.ErrInternalServerError
	}

	name, err := postgres.UploadNewImage(c, buf.Bytes())
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, uploadOut200{name})
}

type (
	postOut200 struct {
		Image string `json:"image"`
	}
)

// @Summary Updates information about user.
// @Security Bearer
// @Accept mpfd
// @Produce json
// @Param image formData file true "image file"
// @Success 200 {object} postOut200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /image/post [post]
func post(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil || file.Filename == "" {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "image is not attached",
			Internal: err,
		}
	}

	src, err := file.Open()
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer closable.SafeClose(src)

	scn := io.Reader(src)
	img, err := imaging.Decode(scn, imaging.AutoOrientation(true))
	if err != nil {
		return echo.ErrInternalServerError
	}

	img = imaging.Fill(img, 600, 600, imaging.Center, imaging.CatmullRom)

	buf := &bytes.Buffer{}
	if err = imaging.Encode(buf, img, imaging.PNG); err != nil {
		return echo.ErrInternalServerError
	}

	name, err := postgres.UploadNewImage(c, buf.Bytes())
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, uploadOut200{name})
}
