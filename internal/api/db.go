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

func (e EndPoint) endpoint() string {
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

func (f *FirebaseClient) Collection(endPoint EndPoint) ([]int, error) {
	ref := f.NewRef(fmt.Sprintf("%s/%s", Version, endPoint.endpoint()))
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
