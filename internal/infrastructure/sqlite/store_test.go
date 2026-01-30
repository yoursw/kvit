package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStore_SetGet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	if err := store.Set(ctx, "servers", "127.0.0.1"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := store.Get(ctx, "servers")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "127.0.0.1" {
		t.Errorf("Get servers = %q, want 127.0.0.1", got)
	}

	if err := store.Set(ctx, "servers/personal", "192.168.1.1"); err != nil {
		t.Fatalf("Set subkey: %v", err)
	}
	got, err = store.Get(ctx, "servers/personal")
	if err != nil {
		t.Fatalf("Get subkey: %v", err)
	}
	if got != "192.168.1.1" {
		t.Errorf("Get servers/personal = %q, want 192.168.1.1", got)
	}

	got, err = store.Get(ctx, "missing")
	if err != nil {
		t.Fatalf("Get missing: %v", err)
	}
	if got != "" {
		t.Errorf("Get missing = %q, want empty", got)
	}
}

func TestStore_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	store.Set(ctx, "k", "v1")
	store.Set(ctx, "k", "v2")
	got, _ := store.Get(ctx, "k")
	if got != "v2" {
		t.Errorf("after overwrite Get = %q, want v2", got)
	}
}

func TestStore_ListKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	// Empty store
	keys, err := store.ListKeys(ctx)
	if err != nil {
		t.Fatalf("ListKeys on empty: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("ListKeys on empty = %v, want empty", keys)
	}

	// Add some keys
	if err := store.Set(ctx, "b", "1"); err != nil {
		t.Fatalf("Set b: %v", err)
	}
	if err := store.Set(ctx, "a", "2"); err != nil {
		t.Fatalf("Set a: %v", err)
	}
	if err := store.Set(ctx, "c", "3"); err != nil {
		t.Fatalf("Set c: %v", err)
	}

	keys, err = store.ListKeys(ctx)
	if err != nil {
		t.Fatalf("ListKeys: %v", err)
	}
	want := []string{"a", "b", "c"}
	if len(keys) != len(want) {
		t.Errorf("ListKeys len = %d, want %d", len(keys), len(want))
	}
	for i, key := range keys {
		if key != want[i] {
			t.Errorf("ListKeys[%d] = %q, want %q", i, key, want[i])
		}
	}
}

func TestNewStore_createsDB(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.db")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	store.Close()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("database file was not created at %s", path)
	}
}
