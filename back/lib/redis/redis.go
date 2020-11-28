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

// Nil is a wrapper of redis.Nil.
const Nil = redis.Nil
const contextKey = "redis"

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

// Inject injects `rd` variable in echo context.
func (rd *Redis) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextKey, rd)
			return next(c)
		}
	}
}

// Store some value in Redis under unique key.
func StoreWithNewToken(c echo.Context, val string, exp time.Duration) (key string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	if tmp, err := uuid.NewRandom(); err != nil {
		return "", err
	} else {
		key = tmp.String()
	}
	for i := 0; i < 3; i++ {
		err = rd.Tokens.SetNX(ctx, key, val, exp).Err()
		if err == nil {
			break
		}
	}
	return
}

// Delete some key from Redis.
func DeleteToken(c echo.Context, key string) error {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Del(ctx, key).Err()
}

func GetValue(c echo.Context, key string) (val string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Get(ctx, key).Result()
}

// Get value and delete key from Redis.
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

// Returns true if token exists in Redis.
func CheckTokenExists(c echo.Context, key string) (bool, error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	val, err := rd.Tokens.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
