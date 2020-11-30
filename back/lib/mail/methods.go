package mail

import (
	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
)

func SendOneTimePassword(c echo.Context, email, password string) error {
	subject := "Login one-time password"
	message := hermes.Email{
		Body: hermes.Body{
			Title: "Hello!",
			Intros: []string{
				"To login at KeepTheMoment.ru please enter the one-time password below.",
				"It will expire in two hours.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Password", Value: password},
			},
		},
	}
	return send(c, email, subject, message, nil)
}
