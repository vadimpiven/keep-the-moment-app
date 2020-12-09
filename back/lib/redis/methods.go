package redis

import (
	"fmt"
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
		err = rd.tokens.SetNX(ctx, key, val, exp).Err()
		if err == nil {
			break
		}
	}
	return
}

func DeleteToken(c echo.Context, key string) (err error) {
	rd, ctx := extract(c)

	return rd.tokens.Del(ctx, key).Err()
}

func GetValueByToken(c echo.Context, key string) (val string, err error) {
	rd := c.Get(contextKey).(*Redis)
	ctx := c.Request().Context()

	return rd.tokens.Get(ctx, key).Result()
}

func GetValueAndDeleteToken(c echo.Context, key string) (val string, err error) {
	rd, ctx := extract(c)

	var tmp *redis.StringCmd
	_, err = rd.tokens.Pipelined(ctx, func(p redis.Pipeliner) error {
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

	return rd.tokens.Expire(ctx, key, exp).Result()
}

func StoreUserCoordsByEmail(c echo.Context, email string, lat, lon float64) (err error) {
	rd, ctx := extract(c)

	res := fmt.Sprintf("%g;%g", lat, lon)
	return rd.coords.Set(ctx, email, res, 0).Err()
}

func GetUserCoordsByEmail(c echo.Context, email string) (lat, lon float64, err error) {
	rd, ctx := extract(c)

	res, err := rd.coords.Get(ctx, email).Result()
	if err != nil {
		return
	}
	_, err = fmt.Sscanf(res, "%g;%g", &lat, &lon)
	return
}
