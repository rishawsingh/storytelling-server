package firebaseProvider

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"os"
	"story-time-server/utils"
)

var FirebaseClient *auth.Client

func InitFirebaseClient() {
	// get firebase key
	fireKey := os.Getenv("FIREBASE_KEY")

	// decode firebase key
	firebaseKey, err := utils.DecodeB64(fireKey)
	if err != nil {
		logrus.Fatalf("unable to get firebase key with error - %v", err)
	}

	opt := option.WithCredentialsJSON(firebaseKey)

	// initiate firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatal("Firebase load error while creating the app", err)
	}

	// initiate firebase auth client
	client, err := app.Auth(context.Background())
	if err != nil {
		logrus.Fatal("Firebase load error", err)
	}

	FirebaseClient = client
	return
}
