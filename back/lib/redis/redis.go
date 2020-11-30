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
	Redis struct {
		tokens *redis.Client
	}
)

func (rd *Redis) Close() error {
	err := rd.tokens.Close()
	return err
}

// New returns new instance of Redis object.
func New(c Config) *Redis {
	tokens := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       0,
	})
	if _, err := tokens.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	fmt.Printf("⇨ redis connection established on [%s]:%d\n", c.Host, c.Port)
	return &Redis{tokens}
}

const contextKey = "__redis__"

// Inject injects Redis in echo context.
func (rd *Redis) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextKey, rd)
			return next(c)
		}
	}
}

func extract(c echo.Context) (rd *Redis, ctx context.Context) {
	rd = c.Get(contextKey).(*Redis)
	ctx = c.Request().Context()
	return
}
