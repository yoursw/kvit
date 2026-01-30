package application

import (
	"context"
	"reflect"
	"testing"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

func TestGetValues_singular(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}
	store.Set(ctx, "config", "value")

	values, err := GetValues(ctx, store, "config", "")
	if err != nil {
		t.Fatalf("GetValues: %v", err)
	}
	expected := []string{"value"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetValues = %v, want %v", values, expected)
	}
}

func TestGetValues_plural_list(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}
	store.Set(ctx, "servers/:0", "127.0.0.1")
	store.Set(ctx, "servers/:1", "192.168.1.1")

	values, err := GetValues(ctx, store, "servers", "")
	if err != nil {
		t.Fatalf("GetValues: %v", err)
	}
	expected := []string{"127.0.0.1", "192.168.1.1"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetValues = %v, want %v", values, expected)
	}
}

func TestGetValues_plural_sublist(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}
	store.Set(ctx, "servers/personal/:0", "home")
	store.Set(ctx, "servers/personal/:1", "work")

	values, err := GetValues(ctx, store, "servers", "personal")
	if err != nil {
		t.Fatalf("GetValues: %v", err)
	}
	expected := []string{"home", "work"}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("GetValues = %v, want %v", values, expected)
	}
}

func TestGetValues_empty(t *testing.T) {
	ctx := context.Background()
	store := &mapStore{m: make(map[string]string)}

	values, err := GetValues(ctx, store, "missing", "")
	if err != nil {
		t.Fatalf("GetValues: %v", err)
	}
	if len(values) != 0 {
		t.Errorf("GetValues empty = %v, want empty", values)
	}
}