package store

import (
	"github.com/CaninoDev/go-hackernews/internal/api"
	"time"
)

type Item struct {
	api.Post
	lastRead time.Time
}

func (i Item) SetReadStamp() {
	i.lastRead = time.Now()
}

func (i Item) GetReadStamp() time.Time {
	return i.lastRead
}
