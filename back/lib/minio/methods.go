package minio

import (
	"bytes"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

const bucketName = "img"

// Methods documented here: https://docs.min.io/docs/golang-client-api-reference

func UploadImage(c echo.Context, img []byte, name string) error {
	mn, ctx := extract(c)

	_, err := mn.PutObject(ctx, bucketName, name, bytes.NewReader(img), int64(len(img)), minio.PutObjectOptions{
		ContentType: "binary/octet-stream",
	})
	return err
}

func (mn *Minio) DeleteImage(ctx context.Context, name string) error {
	return mn.client.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{})
}
