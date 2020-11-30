// This package is a redis.Client wrapper which assists with accessing redis.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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
		Tokens   *redis.Client
		Hashtags *redis.Client
	}
)

func (rd *Redis) Close() error {
	err := rd.Tokens.Close()
	if err != nil {
		return err
	}
	return rd.Hashtags.Close()
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

	hashtags := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       1,
	})
	if _, err := hashtags.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	fmt.Printf("⇨ redis connection established on [%s]:%d\n", c.Host, c.Port)
	return &Redis{tokens, hashtags}
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

// Methods documented here: https://redis.uptrace.dev/

// Store some value in redis under unique key.
func StoreWithNewToken(c echo.Context, val string, exp time.Duration) (key string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	for i := 0; i < 3; i++ {
		if tmp, err := uuid.NewRandom(); err != nil {
			return "", err
		} else {
			key = tmp.String()
		}
		err = rd.Tokens.SetNX(ctx, key, val, exp).Err()
		if err == nil {
			break
		}
	}
	return
}

// Delete some key from redis.
func DeleteToken(c echo.Context, key string) (err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Del(ctx, key).Err()
}

// Get value by key.
func GetValue(c echo.Context, key string) (val string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Get(ctx, key).Result()
}

// Get value and delete key from redis.
func GetValueAndDeleteToken(c echo.Context, key string) (val string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	var tmp *redis.StringCmd
	_, err = rd.Tokens.Pipelined(ctx, func(p redis.Pipeliner) error {
		tmp = p.Get(ctx, key)
		p.Del(ctx, key)
		return nil
	})
	if err != nil {
		return
	}
	return tmp.Val(), nil
}

// Returns true if token exists in redis.
func CheckTokenExistsAndProlong(c echo.Context, key string, exp time.Duration) (val bool, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Expire(ctx, key, exp).Result()
}
