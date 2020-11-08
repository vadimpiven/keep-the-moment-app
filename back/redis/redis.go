// This package is a redis.Client wrapper which assists with accessing redis.
package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type (
	// Config structure contains configurable options of this package.
	Config struct {
		Host     string
		Port     int
		Password string
	}
	// Redis is a redis.Client wrapper.
	Redis redis.Client
)

var ctx = context.Background()

// Nil is a wrapper of redis.Nil.
const Nil = redis.Nil

// New returns new instance of Redis object.
func New(c Config) *Redis {
	rd := (*Redis)(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       0,
	}))
	if _, err := rd.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	fmt.Printf("⇨ redis connection established on [%s]:%d\n", c.Host, c.Port)
	return rd
}

// Inject injects `rd` variable in echo context.
func (rd *Redis) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("rd", rd)
			return next(c)
		}
	}
}
