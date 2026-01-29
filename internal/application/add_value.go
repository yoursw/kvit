package application

import (
	"context"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

// AddValue adds a value for the given key (bucket with optional subkey).
// Key is built as bucket or bucket/subkey via domain.KeyPath.
func AddValue(ctx context.Context, store domain.Store, bucket, subkey, value string) error {
	key := domain.KeyPath(bucket, subkey)
	return store.Set(ctx, key, value)
}
