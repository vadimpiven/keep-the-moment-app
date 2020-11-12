// This package implements router group `auth`.
package auth

import (
	"github.com/sethvargo/go-password/password"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/goware/emailx"
	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
	_ "github.com/sethvargo/go-password/password"

	"github.com/FTi130/keep-the-moment-app/back/lib/jwtclaims"
	"github.com/FTi130/keep-the-moment-app/back/lib/mail"
	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
	_ "github.com/FTi130/keep-the-moment-app/back/postgres"
)

// ApplyRoutes applies routes for the router group.
func ApplyRoutes(g *echo.Group, s []byte) {
	auth := g.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/logout", logout, jwtclaims.Middleware(s))
	}
}

type (
	Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

// login implements login procedure.
// @Summary Generates session token for user
// @Accept  json
// @Produce  json
// @Param credentials body Credentials true "Email and Password"
// @Success 200 {object} map[string]interface{}
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/login [post]
func login(c echo.Context) error {
	cr := new(Credentials)
	if err := c.Bind(cr); err != nil || cr.Email == "" || emailx.Validate(cr.Email) != nil {
		return echo.ErrBadRequest
	}
	cr.Email = emailx.Normalize(cr.Email)

	if cr.Password == "" {
		pw, err := password.Generate(20, 4, 0, false, false)
		if err != nil {
			return echo.ErrInternalServerError
		}
		cr.Password = pw

		rd := c.Get("rd").(*redis.Redis)
		err = rd.Set(c.Request().Context(), cr.Email, pw, time.Hour*2).Err()
		if err != nil {
			return echo.ErrInternalServerError
		}

		em := c.Get("em").(*mail.Email)
		message := hermes.Email{
			Body: hermes.Body{
				Title: "Hello!",
				Intros: []string{
					"To login at KeepTheMoment.ru please enter the one-time password below.",
					"It will expire in two hours.",
				},
				Dictionary: []hermes.Entry{
					{Key: "Password", Value: cr.Password},
				},
			},
		}
		err = em.Send(cr.Email, "Login one-time password", message, nil)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusAccepted, echo.Map{
			"mail": cr.Email,
		})
	}

	rd := c.Get("rd").(*redis.Redis)
	pw, err := rd.Get(c.Request().Context(), cr.Email).Result()
	if err == redis.Nil || pw != cr.Password {
		return echo.ErrBadRequest
	}
	if err != nil {
		return echo.ErrInternalServerError
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtclaims.CustomClaims{
		Email: cr.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})
	sc := c.Get("sc").([]byte)
	t, err := token.SignedString(sc)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// logout implements logout procedure.
func logout(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

// register implements registration of new user.
func register(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}
