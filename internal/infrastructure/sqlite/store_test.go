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
