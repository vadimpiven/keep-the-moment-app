// This package assists in sending service emails.
package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/aymerick/douceur/inliner"
	"github.com/jordan-wright/email"
	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
)

type (
	// Config structure contains configurable options of this package.
	Config struct {
		From     string
		Password string
		Host     string
		Port     int
	}
	// Email is an object providing access to sending emails.
	Email struct {
		auth    smtp.Auth
		addr    string
		from    string
		hermes  hermes.Hermes
		tls     *tls.Config
	}
)

// New returns new instance of Email object.
func New(c Config) *Email {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        "команда Ecokraski.shop",
			Link:        "https://ecokraski.shop/",
			Logo:        "https://ecokraski.shop/ecokraski.png",
			Copyright:   fmt.Sprintf("© %d Ecokraski.shop", time.Now().Year()),
			TroubleText: "Если кнопка {ACTION} не нажимается, скопируйте ссылку ниже в адресную строку браузера.",
		},
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.Host,
	}
	auth := smtp.PlainAuth("", c.From, c.Password, c.Host)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return &Email{auth, addr, c.From, h, tlsConf}
}

// Inject injects `em` variable in echo context.
func (em *Email) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("em", em)
			return next(c)
		}
	}
}

// Send is a function which sends mail to the user.
func (em *Email) Send(to, subject string, message hermes.Email, attachments []*email.Attachment) error {
	text, err := em.hermes.GeneratePlainText(message)
	if err != nil {
		return err
	}
	html, err := em.hermes.GenerateHTML(message)
	if err != nil {
		return err
	}
	html, err = inliner.Inline(html)
	if err != nil {
		return err
	}
	e := &email.Email{
		To:          []string{to},
		From:        em.from,
		Subject:     subject,
		Text:        []byte(text),
		HTML:        []byte(html),
		Headers:     textproto.MIMEHeader{},
		Attachments: attachments,
	}
	return e.SendWithTLS(em.addr, em.auth, em.tls)
}