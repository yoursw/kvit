package domain

import (
	"strconv"
	"strings"
)

// Entry represents a single key-value entry in the store.
// Key is the full key path (e.g. "servers" or "servers/personal").
type Entry struct {
	Key   string
	Value string
}

// KeyPath builds a full key from a bucket and optional subkey.
// e.g. KeyPath("servers", "") -> "servers"; KeyPath("servers", "personal") -> "servers/personal".
func KeyPath(bucket, subkey string) string {
	if subkey == "" {
		return bucket
	}
	return bucket + "/" + subkey
}

// IsPluralBucket reports whether the bucket name denotes a list.
// Plural buckets (e.g. "servers") store positional list entries and named entries;
// singular buckets store a single value or keyed subkeys only.
// Future: list and literal values may support hash associations (named entries).
func IsPluralBucket(bucket string) bool {
	return strings.HasSuffix(bucket, "s") && len(bucket) > 1
}

// List layout: any path can be a list. servers/:0, servers/:1; servers/personal/:0, servers/personal/:1.
const listLenSuffix = "/:len"

func ListLengthKey(listPath string) string {
	return listPath + listLenSuffix
}

func ListItemKey(listPath string, index int) string {
	return listPath + "/:" + strconv.Itoa(index)
}
