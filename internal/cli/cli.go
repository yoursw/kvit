package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thatnerdjosh/kvit/internal/application"
	"github.com/thatnerdjosh/kvit/internal/domain"
	"github.com/thatnerdjosh/kvit/internal/infrastructure/sqlite"
)

const usage = `Usage:
  kvit <bucket> add <value>
  kvit <bucket> add <subkey> <value>

Examples:
  kvit servers add 127.0.0.1
  kvit servers add personal 127.0.0.1
`

// Run parses os.Args and executes the appropriate command.
// Store is created from KVIT_DB (default: $XDG_DATA_HOME/kvit/data.db, or ~/.local/share/kvit/data.db).
func Run(args []string) error {
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}

	bucket := args[0]

	if args[1] != "add" {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}

	rest := args[2:]
	if len(rest) < 1 {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}

	var subkey, value string
	if len(rest) == 1 {
		value = rest[0]
	} else {
		subkey, value = rest[0], rest[1]
	}

	dbPath := os.Getenv("KVIT_DB")
	if dbPath == "" {
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, _ := os.UserHomeDir()
			if home == "" {
				home = "."
			}
			dataHome = filepath.Join(home, ".local", "share")
		}
		dbPath = filepath.Join(dataHome, "kvit", "data.db")
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		return err
	}

	store, err := sqlite.NewStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	key, err := application.AddValue(ctx, store, bucket, subkey, value)
	if err != nil {
		return err
	}

	if domain.IsPluralBucket(bucket) {
		if subkey != "" {
			fmt.Printf("appended to %s (%s): %s\n", bucket, subkey, value)
		} else {
			fmt.Printf("appended to %s: %s\n", bucket, value)
		}
	} else {
		fmt.Printf("stored %s = %s\n", key, value)
	}
	return nil
}
