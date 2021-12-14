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

func TestFirebaseClient_Item(t *testing.T) {
	ctx := context.Background()
	id := 5
	fb, err := NewClientWithDefaults(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fb.Item(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFirebaseClient_Collection(t *testing.T) {
	testCtx := context.Background()
	testClient, err := NewClientWithDefaults(testCtx)
	if err != nil {
		t.Fatal(err)
	}

	testEndpoints := []EndPoint{
		Top, Best, NewS, Jobs, Ask, Show,
	}

	for _, endpoint := range testEndpoints {
		_, err := testClient.CollectionIDs(endpoint)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFirebaseClient_MaxItem(t *testing.T) {
	testCtx := context.Background()
	testClient, err := NewClientWithDefaults(testCtx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testClient.MaxItems()
	if err != nil {
		t.Fatal(err)
	}
}

