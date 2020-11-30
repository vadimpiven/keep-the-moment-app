package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Methods documented here: https://redis.uptrace.dev/

func StoreWithNewToken(c echo.Context, val string, exp time.Duration) (key string, err error) {
	rd, ctx := extract(c)

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

func DeleteToken(c echo.Context, key string) (err error) {
	rd, ctx := extract(c)

	return rd.Tokens.Del(ctx, key).Err()
}

func GetValueByToken(c echo.Context, key string) (val string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.Tokens.Get(ctx, key).Result()
}

func GetValueAndDeleteToken(c echo.Context, key string) (val string, err error) {
	rd, ctx := extract(c)

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

func CheckTokenExistsAndProlong(c echo.Context, key string, exp time.Duration) (val bool, err error) {
	rd, ctx := extract(c)

	return rd.Tokens.Expire(ctx, key, exp).Result()
}
