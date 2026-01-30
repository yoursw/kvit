package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/thatnerdjosh/kvit/internal/application"
	"github.com/thatnerdjosh/kvit/internal/infrastructure/sqlite"
)

func main() {
	var port = flag.Int("port", 14250, "port to listen on")
	flag.Parse()

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
		log.Fatalf("mkdir: %v", err)
	}

	store, err := sqlite.NewStore(dbPath)
	if err != nil {
		log.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	http.HandleFunc("/list-keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		keys, err := application.ListKeys(r.Context(), store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
	})

	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/get/")
		key, err := url.PathUnescape(path)
		if err != nil {
			http.Error(w, "invalid key", http.StatusBadRequest)
			return
		}
		var bucket, subkey string
		if strings.Contains(key, "/") {
			parts := strings.SplitN(key, "/", 2)
			bucket = parts[0]
			subkey = parts[1]
		} else {
			bucket = key
		}
		values, err := application.GetValues(r.Context(), store, bucket, subkey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(values)
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("kvitd listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}