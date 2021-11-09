package api

import (
	"context"
	"testing"
)

func TestNewClientWithDefaults(t *testing.T) {
	ctx := context.Background()

	_, err := NewClientWithDefaults(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
