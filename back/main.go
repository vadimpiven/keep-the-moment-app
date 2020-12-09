package main

import (
	"github.com/FTi130/keep-the-moment-app/back/lib/closable"
	"github.com/FTi130/keep-the-moment-app/back/lib/config"
	"github.com/FTi130/keep-the-moment-app/back/lib/flags"
	"github.com/FTi130/keep-the-moment-app/back/lib/mail"
	"github.com/FTi130/keep-the-moment-app/back/lib/minio"
	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
	"github.com/FTi130/keep-the-moment-app/back/lib/redis"
	"github.com/FTi130/keep-the-moment-app/back/lib/watchdog"
	"github.com/FTi130/keep-the-moment-app/back/server"
)

// @title KeepTheMoment REST API
// @version 1.0

// @host keepthemoment.ru
// @BasePath /api
// @schemes https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	// init configs
	f := flags.Read()
	c := config.Read(f.Config)
	s := server.New(f.Server, c.Server)

	// postgres
	db := postgres.New(c.Postgres)
	defer closable.SafeClose(db)
	s.Pre(db.Inject())

	// redis
	rd := redis.New(c.Redis)
	defer closable.SafeClose(rd)
	s.Pre(rd.Inject())

	//minio
	mn := minio.New(c.Minio)
	s.Pre(mn.Inject())

	// mail
	em := mail.New(c.Email)
	s.Pre(em.Inject())

	// watchdog
	watchdog.Watch(db, mn)

	// entry point
	s.Run()
}
