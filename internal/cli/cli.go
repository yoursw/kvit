package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/thatnerdjosh/kvit/internal/application"
	"github.com/thatnerdjosh/kvit/internal/domain"
	"github.com/thatnerdjosh/kvit/internal/infrastructure/sqlite"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Current   string            `json:"current"`
	Contexts  map[string]string `json:"contexts"`
}

func getConfigPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		if home == "" {
			home = "."
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "kvit", "config.yaml")
}

func loadConfig() (*Config, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Contexts: make(map[string]string)}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]string)
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

const usage = `Usage:
  kvit <bucket> add <value>
  kvit <bucket> add <subkey> <value>
  kvit <bucket> get [subkey]
  kvit list-keys
  kvit context add <name> <address>
  kvit context use <name>
  kvit context unset

Examples:
  kvit servers add 127.0.0.1
  kvit servers add personal 127.0.0.1
  kvit servers get
  kvit servers get personal
  kvit config get db
  kvit list-keys
  kvit context add loopback 127.0.0.1
  kvit context use loopback
  kvit context unset
`

// Run parses os.Args and executes the appropriate command.
// Store is created from KVIT_DB (default: $XDG_DATA_HOME/kvit/data.db, or ~/.local/share/kvit/data.db).
func Run(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}

	if args[0] == "list-keys" {
		if len(args) != 1 {
			fmt.Fprint(os.Stderr, usage)
			return nil
		}
		return runListKeys()
	}

	if args[0] == "context" {
		if len(args) < 2 {
			fmt.Fprint(os.Stderr, usage)
			return nil
		}
		sub := args[1]
		switch sub {
		case "add":
			if len(args) != 4 {
				fmt.Fprint(os.Stderr, usage)
				return nil
			}
			name := args[2]
			address := args[3]
			return runContextAdd(name, address)
		case "use":
			if len(args) != 3 {
				fmt.Fprint(os.Stderr, usage)
				return nil
			}
			name := args[2]
			return runContextUse(name)
		case "unset":
			if len(args) != 2 {
				fmt.Fprint(os.Stderr, usage)
				return nil
			}
			return runContextUnset()
		default:
			fmt.Fprint(os.Stderr, usage)
			return nil
		}
	}

	if len(args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}

	bucket := args[0]
	cmd := args[1]

	if cmd == "add" {
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

		dbPath := getDBPath()
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
	} else if cmd == "get" {
		subkey := ""
		if len(args) > 2 {
			subkey = args[2]
		}
		return runGet(bucket, subkey)
	} else {
		fmt.Fprint(os.Stderr, usage)
		return nil
	}
}

func runListKeys() error {
	// Check if remote context is set
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if cfg.Current != "" {
		address, ok := cfg.Contexts[cfg.Current]
		if ok && address != "" {
			fmt.Fprintf(os.Stderr, "Warning: Data is sent unencrypted over the network. Ensure sensitive data is properly encapsulated.\n")
			// Make HTTP request to remote kvitd
			resp, err := http.Get("http://" + address + ":14250/list-keys")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("remote error: %s", resp.Status)
			}
			var keys []string
			if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
				return err
			}
			for _, key := range keys {
				fmt.Println(key)
			}
			return nil
		}
		// If address empty or not found, fall back to local
	}

	// Local execution
	dbPath := getDBPath()
	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		return err
	}

	store, err := sqlite.NewStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	keys, err := application.ListKeys(ctx, store)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Println(key)
	}
	return nil
}

func runGet(bucket, subkey string) error {
	// Check if remote context is set
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	key := domain.KeyPath(bucket, subkey)
	if cfg.Current != "" {
		address, ok := cfg.Contexts[cfg.Current]
		if ok && address != "" {
			fmt.Fprintf(os.Stderr, "Warning: Data is sent unencrypted over the network. Ensure sensitive data is properly encapsulated.\n")
			// Make HTTP request to remote kvitd
			resp, err := http.Get("http://" + address + ":14250/get/" + url.PathEscape(key))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("remote error: %s", resp.Status)
			}
			var values []string
			if err := json.NewDecoder(resp.Body).Decode(&values); err != nil {
				return err
			}
			for _, value := range values {
				fmt.Println(value)
			}
			return nil
		}
		// If address empty or not found, fall back to local
	}

	// Local execution
	dbPath := getDBPath()
	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		return err
	}

	store, err := sqlite.NewStore(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	values, err := application.GetValues(ctx, store, bucket, subkey)
	if err != nil {
		return err
	}

	for _, value := range values {
		fmt.Println(value)
	}
	return nil
}

func runContextAdd(name, address string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.Contexts[name] = address
	return saveConfig(cfg)
}

func runContextUse(name string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if _, ok := cfg.Contexts[name]; !ok {
		return fmt.Errorf("context %q not found", name)
	}
	cfg.Current = name
	return saveConfig(cfg)
}

func runContextUnset() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.Current = ""
	return saveConfig(cfg)
}

func getDBPath() string {
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
	return dbPath
}
