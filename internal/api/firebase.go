package api

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
)

const (
	BaseURI = "https://hacker-news.firebaseio.com/"
	Version = "v0"
)

type FirebaseClient struct {
	ctx context.Context
	*db.Client
}

func NewClientWithDefaults(ctx context.Context) (*FirebaseClient, error) {
	defaultCfg := &firebase.Config{
		DatabaseURL: BaseURI,
	}

	app, err := firebase.NewApp(ctx, defaultCfg, option.WithoutAuthentication())
	if err != nil {
		errorString := fmt.Sprintf("error initializing firebase app: %v", err)
		return &FirebaseClient{}, errors.New(errorString)
	}

	fb, err := app.Database(ctx)
	if err != nil {
		errorString := fmt.Sprintf("error initializing firebase connection: %v", err)
		return &FirebaseClient{}, errors.New(errorString)
	}

	return &FirebaseClient{
		ctx,
		fb,
	}, nil
}


