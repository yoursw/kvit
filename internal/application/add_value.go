package application

import (
	"context"
	"strconv"

	"github.com/thatnerdjosh/kvit/internal/domain"
)

// AddValue adds a value. Plural bucket: append to list at bucket or bucket/subkey.
// servers add 127.0.0.1 → servers/:0; servers add personal 127.0.0.1 → servers/personal/:0.
// Singular: set bucket or bucket/subkey.
func AddValue(ctx context.Context, store domain.Store, bucket, subkey, value string) (effectiveKey string, err error) {
	if domain.IsPluralBucket(bucket) {
		listPath := domain.KeyPath(bucket, subkey) // "servers" or "servers/personal"
		lenKey := domain.ListLengthKey(listPath)
		lenStr, err := store.Get(ctx, lenKey)
		if err != nil {
			return "", err
		}
		n := 0
		if lenStr != "" {
			n, err = strconv.Atoi(lenStr)
			if err != nil {
				return "", err
			}
		}
		itemKey := domain.ListItemKey(listPath, n)
		if err := store.Set(ctx, itemKey, value); err != nil {
			return "", err
		}
		if err := store.Set(ctx, lenKey, strconv.Itoa(n+1)); err != nil {
			return "", err
		}
		return itemKey, nil
	}
	key := domain.KeyPath(bucket, subkey)
	if err := store.Set(ctx, key, value); err != nil {
		return "", err
	}
	return key, nil
}
