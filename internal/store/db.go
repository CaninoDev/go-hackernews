package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"log"
	"sync"
	"time"
)

type Store struct {
	fb            *api.FirebaseClient
	items         map[int]Item
	idCollections map[api.EndPoint][]int
	mutex         sync.RWMutex
}

func New() (*Store, error) {
	ctx := context.Background()
	fb, err := api.NewClientWithDefaults(ctx)
	if err != nil {
		return &Store{}, err
	}

	store := &Store{
		fb:            fb,
		items:         make(map[int]Item),
		idCollections: make(map[api.EndPoint][]int),
	}
	store.cacheCollections()
	return store, nil
}

func (s *Store) Collection(endPointStr string) []int {
	endPoint := api.ToEndPoint(endPointStr)
	return s.idCollections[endPoint]
}

func (s *Store) Item(id int) (Item, error) {
	return s.item(id)
}

func (s *Store) SetItemReadStamp(id int) error {
	if _, ok := s.items[id]; ok {
		s.items[id].SetReadStamp()
		return nil
	} else {
		return errors.New(fmt.Sprintf("Item id: %d not found", id))
	}
}

func (s *Store) GetItemReadStamp(id int) (time.Time, error) {
	if val, ok := s.items[id]; ok {
		return val.GetReadStamp(), nil
	} else {
		return time.Unix(0, 0), errors.New(fmt.Sprintf("Item id: %d not found", id))
	}
}

func (s *Store) Subscribe(ids []int) Subscription {
	sub := &subscription{
		ids:     ids,
		store:   s,
		updates: make(chan Item),
	}
	go sub.loop()
	return sub
}

func (s *Store) CollectionsList() (endPointStrings []string) {
	endPoints := api.AllEndPoints()
	for _, endPoint := range endPoints {
		endPointStrings = append(endPointStrings, endPoint.String())
	}
	return endPointStrings
}

func (s *Store) cacheCollections() {
	for _, endPoint := range api.AllEndPoints() {
		var err error
		s.idCollections[endPoint], err = s.fb.CollectionIDs(endPoint)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Store) item(id int) (Item, error) {
	if val, ok := s.items[id]; ok {
		return val, nil
	} else {
		post, err := s.fb.Item(id)
		if err != nil {
			return Item{post, time.Unix(0, 0)}, err
		}
		s.items[id] = Item{post, time.Unix(0, 0)}
		return s.items[id], nil
	}
}
