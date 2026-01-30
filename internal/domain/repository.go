package domain

import "context"

// Store is the repository interface for key-value persistence.
// Implementations can be SQLite, Redis, or any other backend.
type Store interface {
	// Set stores or overwrites the value for the given key.
	Set(ctx context.Context, key, value string) error
	// Get returns the value for the given key, or empty string if not found.
	Get(ctx context.Context, key string) (string, error)
	// ListKeys returns all keys in the store, sorted.
	ListKeys(ctx context.Context) ([]string, error)
	// Close releases resources. Callers should invoke it when done.
	Close() error
}
