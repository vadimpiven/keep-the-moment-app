package firebase

import (
	"context"

	"firebase.google.com/go/messaging"
)

func (fb *Firebase) Notify(ctx context.Context, token, title, message string) (tokenValid bool, err error) {
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	}
	_, err = fb.messaging.Send(ctx, msg)
	tokenValid = messaging.IsRegistrationTokenNotRegistered(err)
	return
}
