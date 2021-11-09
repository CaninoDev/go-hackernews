package api

import (
	"sync"
)

type Subscription interface {
	Updates() <-chan Post
	Close() error
}


type subscription struct {
	postIDs []int
	db *FirebaseClient
	updates chan Post
	closing chan chan error
}

func (s *subscription) Updates() <-chan Post {
	return s.updates
}

func (s *subscription) Close() error {
	errCh := make(chan error)
	s.closing <-errCh
	return <-errCh
}

func (s *subscription) loop() {
	defer close(s.updates)
	var err error
	var wg sync.WaitGroup
	var post Post

	reset := func() {
		if len(s.postIDs) == 0 {
			<-s.closing
			return
		}
	}
	reset()

	for _, id := range s.postIDs {
		wg.Add(1)

		post, err = s.db.Item(id)
		select {
		case errCh := <-s.closing:
			errCh <-err
			close(s.updates)
			return
		default:
			s.updates <-post
		}
		wg.Done()
	}

}

func (f *FirebaseClient) Subscribe(IDs []int) Subscription {
	s := &subscription{
		postIDs: IDs,
		db: f,
		updates: make(chan Post),
	}
	go s.loop()
	return s
}