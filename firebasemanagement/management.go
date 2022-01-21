package firebasemanagement

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App

func init() {
	opt := option.WithCredentialsFile("firebase/fitnesstime-2345e-firebase-adminsdk-dai66-f54354e14f.json")

	var err error

	App, err = firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		fmt.Print(err)
	}
}

// SendMessage ...
func SendMessage(title string, messageContent string, fcmToken string) (string, error) {

	ctx := context.Background()

	messagingClient, err := App.Messaging(ctx)

	notification := &messaging.Notification{
		Title: title,
		Body:  messageContent,
	}

	message := &messaging.Message{
		Token:        fcmToken,
		Notification: notification,
	}

	res, err := messagingClient.Send(ctx, message)

	return res, err
}
