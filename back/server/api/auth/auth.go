// This package implements router group `auth`.
package auth

import (
	"net/http"
	"time"

	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"

	"github.com/goware/emailx"
	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"

	"github.com/FTi130/keep-the-moment-app/back/lib/keyauth"
	"github.com/FTi130/keep-the-moment-app/back/lib/mail"
	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group) {
	auth := g.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/logout", logout, keyauth.Middleware())
	}
}

type (
	loginIn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	loginOut200 struct {
		Token string `json:"token"`
	}
	loginOut202 struct {
		Email string `json:"email"`
	}
)

// Generates session token for user.
// @Summary Generates session token for user.
// @Accept json
// @Produce json
// @Param credentials body loginIn true "email and password"
// @Success 200 {object} loginOut200
// @Success 202 {object} loginOut202
// @Failure 400,500 {object} httputil.HTTPError
// @Router /auth/login [post]
func login(c echo.Context) error {
	cr := new(loginIn)
	err := c.Bind(cr)
	if err != nil {
		return echo.ErrBadRequest
	}

	if cr.Password == "" {
		if cr.Email == "" || emailx.Validate(cr.Email) != nil {
			return echo.ErrBadRequest
		}
		cr.Email = emailx.Normalize(cr.Email)

		token, err := redis.StoreWithNewToken(c, cr.Email, 2*time.Hour)
		if err != nil {
			return echo.ErrInternalServerError
		}

		message := hermes.Email{
			Body: hermes.Body{
				Title: "Hello!",
				Intros: []string{
					"To login at KeepTheMoment.ru please enter the one-time password below.",
					"It will expire in two hours.",
				},
				Dictionary: []hermes.Entry{
					{Key: "Password", Value: token},
				},
			},
		}
		err = mail.Send(c, cr.Email, "Login one-time password", message, nil)

		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusAccepted, loginOut202{
			Email: cr.Email,
		})
	}

	val, err := redis.GetValueAndDeleteToken(c, cr.Password)
	if err != nil {
		return echo.ErrInternalServerError
	} else if val != cr.Email {
		return echo.ErrBadRequest
	}

	token, err := redis.StoreWithNewToken(c, cr.Email, 72*time.Hour)
	if err != nil {
		return echo.ErrInternalServerError
	}

	err = postgres.RegisterIfNewUser(c, cr.Email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, loginOut200{
		Token: token,
	})
}

// Expires session token.
// @Summary Expires session token.
// @Produce json
// @Success 200
// @Failure 400,401,500 {object} httputil.HTTPError
// @Router /auth/logout [post]
func logout(c echo.Context) error {
	if err := keyauth.ExpireToken(c); err != nil {
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}
