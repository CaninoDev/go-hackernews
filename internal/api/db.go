package api

import (
	"fmt"
)

type EndPoint int

const (
	Top EndPoint = iota
	Best
	NewS
	Ask
	Show
	Jobs
)

var toEndPoint = map[string]EndPoint{
	"Top": Top,
	"Best": Best,
	"New": NewS,
	"Ask": Ask,
	"Show": Show,
	"Jobs": Jobs,
}

func (e EndPoint) endpointURL() string {
	return [...]string{"topstories", "beststories", "newstories", "askstories", "showstories", "jobstories"}[e]
}

func (e EndPoint) String() string {
	return [...]string{"Top", "Best", "New", "Ask", "Show", "Jobs"}[e]
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

func (f *FirebaseClient) Collection(endPointStr string) ([]int, error) {
	endPoint := toEndPoint[endPointStr]
	ref := f.NewRef(fmt.Sprintf("%s/%s", Version, endPoint.endpointURL()))
	var stories []int
	err := ref.Get(f.ctx, &stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

func (f *FirebaseClient) MaxItem() (int, error) {
	ref := f.NewRef(fmt.Sprintf("%s/maxitem", Version))
	var maxItem int
	err := ref.Get(f.ctx, &maxItem)
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

func (f *FirebaseClient) EndPoints() []string {
	return []string{
		Top.String(), Best.String(), NewS.String(), Ask.String(), Show.String(), Jobs.String(),
	}
}
