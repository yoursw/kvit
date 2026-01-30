package application

import (
	"context"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

// ListKeys returns all keys in the store.
func ListKeys(ctx context.Context, store domain.Store) ([]string, error) {
	return store.ListKeys(ctx)
}