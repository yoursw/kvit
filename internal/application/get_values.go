package application

import (
	"context"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

// GetValues returns the values for the given bucket and optional subkey.
// For plural buckets, returns the list of values.
// For singular buckets, returns a single value.
func GetValues(ctx context.Context, store domain.Store, bucket, subkey string) ([]string, error) {
	key := domain.KeyPath(bucket, subkey)
	if domain.IsPluralBucket(bucket) {
		var values []string
		i := 0
		for {
			itemKey := domain.ListItemKey(key, i)
			value, err := store.Get(ctx, itemKey)
			if err != nil {
				return nil, err
			}
			if value == "" {
				break
			}
			values = append(values, value)
			i++
		}
		return values, nil
	} else {
		value, err := store.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		if value == "" {
			return []string{}, nil
		}
		return []string{value}, nil
	}
}