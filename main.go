package main

import (
	"context"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/CaninoDev/go-hackernews/internal/ui"
	"log"
)

func main() {
	ctx := context.Background()

	db, err := api.NewClientWithDefaults(ctx)
	if err != nil {
		log.Fatal(err)
	}

	displayEngine := ui.Init(db)

	if err := displayEngine.Run(); err != nil {
		log.Fatal(err)
	}
}
