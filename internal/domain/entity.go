package domain

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
