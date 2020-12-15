// This package is a Firebase messaging.Client wrapper which assists with sending push notifications.
package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type (
	// Config structure contains configurable options of this package.
	Config struct {
		APIKey    string
		ProjectID string
	}
	// Minio is a minio.Client wrapper.
	Firebase struct {
		messaging *messaging.Client
	}
)

// New returns new instance of Firebase object.
func New(c Config) (fb *Firebase) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: c.ProjectID,
	}, option.WithAPIKey(c.APIKey))
	if err != nil {
		panic(err)
	}
	msg, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("⇨ Firebase connection established")
	return &Firebase{msg}
}
