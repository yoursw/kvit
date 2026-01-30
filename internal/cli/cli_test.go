package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thatnerdjosh/kvit/internal/infrastructure/sqlite"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	// Test loading non-existent file
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Current != "" {
		t.Errorf("expected empty current, got %q", cfg.Current)
	}
	if len(cfg.Contexts) != 0 {
		t.Errorf("expected empty contexts, got %v", cfg.Contexts)
	}
}

func TestSaveLoadConfig(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	cfg := &Config{
		Current: "test",
		Contexts: map[string]string{
			"test": "127.0.0.1",
		},
	}

	if err := saveConfig(cfg); err != nil {
		t.Fatalf("saveConfig: %v", err)
	}

	loaded, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}

	if loaded.Current != "test" {
		t.Errorf("expected current 'test', got %q", loaded.Current)
	}
	if loaded.Contexts["test"] != "127.0.0.1" {
		t.Errorf("expected contexts[test] '127.0.0.1', got %q", loaded.Contexts["test"])
	}
}

func TestRunContextAdd(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	if err := runContextAdd("loopback", "127.0.0.1"); err != nil {
		t.Fatalf("runContextAdd: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Contexts["loopback"] != "127.0.0.1" {
		t.Errorf("expected loopback '127.0.0.1', got %q", cfg.Contexts["loopback"])
	}
}

func TestRunContextUse(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	// Add first
	if err := runContextAdd("loopback", "127.0.0.1"); err != nil {
		t.Fatalf("runContextAdd: %v", err)
	}

	// Use
	if err := runContextUse("loopback"); err != nil {
		t.Fatalf("runContextUse: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Current != "loopback" {
		t.Errorf("expected current 'loopback', got %q", cfg.Current)
	}

	// Use non-existent
	if err := runContextUse("nonexist"); err == nil {
		t.Error("expected error for non-existent context")
	}
}

func TestRunContextUnset(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	// Add and use
	if err := runContextAdd("loopback", "127.0.0.1"); err != nil {
		t.Fatalf("runContextAdd: %v", err)
	}
	if err := runContextUse("loopback"); err != nil {
		t.Fatalf("runContextUse: %v", err)
	}

	// Unset
	if err := runContextUnset(); err != nil {
		t.Fatalf("runContextUnset: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Current != "" {
		t.Errorf("expected empty current, got %q", cfg.Current)
	}
}

func TestRun_invalidArgs(t *testing.T) {
	// Test Run with invalid args
	err := Run([]string{})
	if err != nil {
		t.Errorf("Run with empty args should not return error, got %v", err)
	}

	err = Run([]string{"invalid"})
	if err != nil {
		t.Errorf("Run with invalid command should not return error, got %v", err)
	}
}

func TestRunListKeys_local(t *testing.T) {
	dir := t.TempDir()
	oldDB := os.Getenv("KVIT_DB")
	oldConfig := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("KVIT_DB", filepath.Join(dir, "test.db"))
	os.Setenv("XDG_CONFIG_HOME", dir)
	defer func() {
		os.Setenv("KVIT_DB", oldDB)
		os.Setenv("XDG_CONFIG_HOME", oldConfig)
	}()

	// Setup some data
	store, err := sqlite.NewStore(os.Getenv("KVIT_DB"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	ctx := context.Background()
	store.Set(ctx, "test", "value")
	store.Close()

	// Capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = runListKeys()

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("runListKeys: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test") {
		t.Errorf("expected output to contain 'test', got %q", output)
	}
}

func TestGetConfigPath(t *testing.T) {
	old := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", old)

	// With XDG_CONFIG_HOME set
	os.Setenv("XDG_CONFIG_HOME", "/tmp")
	path := getConfigPath()
	expected := filepath.Join("/tmp", "kvit", "config.yaml")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}

	// Without
	os.Setenv("XDG_CONFIG_HOME", "")
	path = getConfigPath()
	home, _ := os.UserHomeDir()
	expected = filepath.Join(home, ".config", "kvit", "config.yaml")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}