package application

import (
	"context"
	"testing"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

type fakeStore struct {
	key, value string
}

func (f *fakeStore) Set(ctx context.Context, key, value string) error {
	f.key, f.value = key, value
	return nil
}

func (f *fakeStore) Get(ctx context.Context, key string) (string, error) {
	if key == f.key {
		return f.value, nil
	}
	return "", nil
}

func (f *fakeStore) Close() error { return nil }

func TestAddValue(t *testing.T) {
	ctx := context.Background()
	store := &fakeStore{}

	err := AddValue(ctx, store, "servers", "", "127.0.0.1")
	if err != nil {
		t.Fatalf("AddValue: %v", err)
	}
	if store.key != "servers" || store.value != "127.0.0.1" {
		t.Errorf("store: key=%q value=%q, want servers / 127.0.0.1", store.key, store.value)
	}

	err = AddValue(ctx, store, "servers", "personal", "192.168.1.1")
	if err != nil {
		t.Fatalf("AddValue with subkey: %v", err)
	}
	wantKey := domain.KeyPath("servers", "personal")
	if store.key != wantKey || store.value != "192.168.1.1" {
		t.Errorf("store: key=%q value=%q, want %q / 192.168.1.1", store.key, store.value, wantKey)
	}
}
