package mail

import (
	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
)

func SendOneTimePassword(c echo.Context, email, password string) error {
	subject := "One-time verification code"
	message := hermes.Email{
		Body: hermes.Body{
			Title: "Hello!",
			Intros: []string{
				"To login at KeepTheMoment.ru please enter the one-time verification code below.",
				"It will expire in two hours.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Verification code", Value: password},
			},
		},
	}
	return send(c, email, subject, message, nil)
}
