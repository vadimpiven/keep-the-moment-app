package main

import (
	"github.com/FTi130/keep-the-moment-app/back/lib/closable"
	"github.com/FTi130/keep-the-moment-app/back/lib/config"
	"github.com/FTi130/keep-the-moment-app/back/lib/flags"
	"github.com/FTi130/keep-the-moment-app/back/lib/mail"
	"github.com/FTi130/keep-the-moment-app/back/postgres"
	"github.com/FTi130/keep-the-moment-app/back/redis"
	"github.com/FTi130/keep-the-moment-app/back/server"
)

// @title KeepTheMoment REST API
// @version 1.0

// @host keepthemoment.ru
// @BasePath /api
func main() {
	// init configs
	f := flags.Read()
	c := config.Read(f.Config)
	s := server.New(f.Server, c.Server)

	// postgres
	db := postgres.New(c.Postgres)
	defer closable.SafeClose(db)
	s.Use(db.Inject())

	// redis
	rd := redis.New(c.Redis)
	defer closable.SafeClose(rd)
	s.Use(rd.Inject())

	// mail
	em := mail.New(c.Email)
	s.Use(em.Inject())

	// entry point
	s.Run()
}
