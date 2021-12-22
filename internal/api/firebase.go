package api

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	fb "firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
)

const (
	BaseURI = "https://hacker-news.firebaseio.com/"
	Version = "v0"
)

//go:generate stringer -type=EndPoint
type EndPoint int

const (
	Top EndPoint = iota
	New
	Best
	Ask
	Show
	Jobs
)

var endPointURL = map[EndPoint]string{
	Top:  "topstories",
	New:  "newstories",
	Best: "beststories",
	Ask:  "askstories",
	Show: "showstories",
	Jobs: "jobstories",
}

var toEndPoint = map[string]EndPoint{
	"Top":  Top,
	"New":  New,
	"Best": Best,
	"Ask":  Ask,
	"Show": Show,
	"Jobs": Jobs,
}

type FirebaseClient struct {
	ctx context.Context
	*fb.Client
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

func (f *FirebaseClient) Item(id int) (Post, error) {
	ref := f.NewRef(fmt.Sprintf("%s/item/%d", Version, id))
	var item post
	err := ref.Get(f.ctx, &item)
	if err != nil {
		return nil, err
	}
	return item, err
}

func (f *FirebaseClient) CollectionIDs(endPoint EndPoint) ([]int, error) {
	ref := f.NewRef(fmt.Sprintf("%s/%s", Version, endPointURL[endPoint]))
	var ids []int
	if err := ref.Get(f.ctx, &ids); err != nil {
		return ids, err
	}
	return ids, nil
}

func (f *FirebaseClient) MaxItems() (int, error) {
	ref := f.NewRef(fmt.Sprintf("%s/maxitem", Version))
	var maxItem int
	err := ref.Get(f.ctx, &maxItem)
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

func AllEndPoints() []EndPoint {
	endpoints := make([]EndPoint, Jobs)
	for i := 0; i < int(Jobs); i++ {
		endpoints[i] = EndPoint(i)
	}
	return endpoints
}

func ToEndPoint(endPointStr string) EndPoint {
	return toEndPoint[endPointStr]
}
