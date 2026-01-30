package application

import (
	"context"
	"reflect"
	"testing"
)

func TestListKeys(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}

	// Empty
	keys, err := ListKeys(ctx, store)
	if err != nil {
		t.Fatalf("ListKeys empty: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("ListKeys empty len = %d, want 0", len(keys))
	}

	// Add some
	store.Set(ctx, "z", "1")
	store.Set(ctx, "a", "2")
	store.Set(ctx, "m", "3")

	keys, err = ListKeys(ctx, store)
	if err != nil {
		t.Fatalf("ListKeys: %v", err)
	}
	want := []string{"a", "m", "z"}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("ListKeys = %v, want %v", keys, want)
	}
}