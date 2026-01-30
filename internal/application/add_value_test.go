package application

import (
	"context"
	"sort"
	"testing"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

// mapStore implements domain.Store with an in-memory map (supports multiple keys for list behavior).
type mapStore struct {
	m map[string]string
}

func (s *mapStore) Set(ctx context.Context, key, value string) error {
	if s.m == nil {
		s.m = make(map[string]string)
	}
	s.m[key] = value
	return nil
}

func (s *mapStore) Get(ctx context.Context, key string) (string, error) {
	return s.m[key], nil
}

func (s *mapStore) ListKeys(ctx context.Context) ([]string, error) {
	var keys []string
	for k := range s.m {
		keys = append(keys, k)
	}
	// Sort for consistency
	sort.Strings(keys)
	return keys, nil
}

func (s *mapStore) Close() error { return nil }

func TestAddValue_singular(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}

	key, err := AddValue(ctx, store, "config", "", "value")
	if err != nil {
		t.Fatalf("AddValue: %v", err)
	}
	if key != "config" {
		t.Errorf("key = %q, want config", key)
	}
	if store.m["config"] != "value" {
		t.Errorf("config = %q, want value", store.m["config"])
	}

	key, err = AddValue(ctx, store, "config", "db", "postgres")
	if err != nil {
		t.Fatalf("AddValue subkey: %v", err)
	}
	if key != "config/db" {
		t.Errorf("key = %q, want config/db", key)
	}
	if store.m["config/db"] != "postgres" {
		t.Errorf("config/db = %q, want postgres", store.m["config/db"])
	}
}

func TestAddValue_plural_listAppend(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}

	key, err := AddValue(ctx, store, "servers", "", "127.0.0.1")
	if err != nil {
		t.Fatalf("AddValue: %v", err)
	}
	if key != "servers/:0" {
		t.Errorf("key = %q, want servers/:0", key)
	}
	if store.m["servers/:0"] != "127.0.0.1" {
		t.Errorf("servers/:0 = %q, want 127.0.0.1", store.m["servers/:0"])
	}
	if store.m[domain.ListLengthKey("servers")] != "1" {
		t.Errorf("len = %q, want 1", store.m[domain.ListLengthKey("servers")])
	}

	key, err = AddValue(ctx, store, "servers", "", "192.168.1.1")
	if err != nil {
		t.Fatalf("AddValue second: %v", err)
	}
	if key != "servers/:1" {
		t.Errorf("key = %q, want servers/:1", key)
	}
	if store.m["servers/:1"] != "192.168.1.1" {
		t.Errorf("servers/:1 = %q, want 192.168.1.1", store.m["servers/:1"])
	}
	if store.m[domain.ListLengthKey("servers")] != "2" {
		t.Errorf("len = %q, want 2", store.m[domain.ListLengthKey("servers")])
	}
}

func TestAddValue_plural_namedEntry(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}

	key, err := AddValue(ctx, store, "servers", "personal", "192.168.1.1")
	if err != nil {
		t.Fatalf("AddValue: %v", err)
	}
	// List underneath: servers/personal/:0
	if key != "servers/personal/:0" {
		t.Errorf("key = %q, want servers/personal/:0", key)
	}
	if store.m["servers/personal/:0"] != "192.168.1.1" {
		t.Errorf("servers/personal/:0 = %q, want 192.168.1.1", store.m["servers/personal/:0"])
	}
	if store.m[domain.ListLengthKey("servers/personal")] != "1" {
		t.Errorf("servers/personal/:len = %q, want 1", store.m[domain.ListLengthKey("servers/personal")])
	}
}
