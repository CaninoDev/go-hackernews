package main

import (
	"github.com/CaninoDev/go-hackernews/internal/ui"
	"log"
)

func main() {

	displayEngine, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := displayEngine.Start(); err != nil {
		log.Fatal(err)
	}
}
