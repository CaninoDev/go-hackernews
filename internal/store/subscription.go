package store

import (
	"sync"
)

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

type subscription struct {
	ids     []int
	store   *Store
	updates chan Item
	closing chan chan error
}

func (s *subscription) Updates() <-chan Item {
	return s.updates
}

func (s *subscription) Close() error {
	errCh := make(chan error)
	s.closing <- errCh
	return <-errCh
}

func (s *subscription) loop() {
	defer close(s.updates)

	var wg sync.WaitGroup

	reset := func() {
		if len(s.ids) == 0 {
			<-s.closing
			return
		}
	}
	reset()

	for _, id := range s.ids {
		var err error
		var post Item

		wg.Add(1)

		post, err = s.store.Item(id)

		select {
		case errCh := <-s.closing:
			errCh <- err
			close(s.updates)
			return
		default:
			s.updates <- post
		}
		wg.Done()
	}
}
