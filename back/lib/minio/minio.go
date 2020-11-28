// This package is a minio.Client wrapper which assists with accessing MinIO.
package minio

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type (
	// Config structure contains configurable options of this package.
	Config struct {
		Host      string
		Port      int
		AccessKey string
		SecretKey string
	}
	// Minio is a minio.Client wrapper.
	Minio minio.Client
)

// New returns new instance of Redis object.
func New(c Config) *Minio {
	mn, err := minio.New(fmt.Sprintf("%s:%d", c.Host, c.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	exists, err := mn.BucketExists(context.Background(), "img")
	if err != nil || !exists {
		panic(err)
	}
	fmt.Printf("⇨ MinIO connection established on [%s]:%d\n", c.Host, c.Port)
	return (*Minio)(mn)
}

// Inject injects `mn` variable in echo context.
func (mn *Minio) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mn", mn)
			return next(c)
		}
	}
}
