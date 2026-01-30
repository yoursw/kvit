package domain

import "testing"

func TestKeyPath(t *testing.T) {
	tests := []struct {
		bucket, subkey, want string
	}{
		{"servers", "", "servers"},
		{"servers", "personal", "servers/personal"},
		{"config", "db", "config/db"},
	}
	for _, tt := range tests {
		got := KeyPath(tt.bucket, tt.subkey)
		if got != tt.want {
			t.Errorf("KeyPath(%q, %q) = %q, want %q", tt.bucket, tt.subkey, got, tt.want)
		}
	}
}

func TestIsPluralBucket(t *testing.T) {
	plural := []string{"servers", "hosts", "keys"}
	for _, b := range plural {
		if !IsPluralBucket(b) {
			t.Errorf("IsPluralBucket(%q) = false, want true", b)
		}
	}
	singular := []string{"server", "config", "key"}
	for _, b := range singular {
		if IsPluralBucket(b) {
			t.Errorf("IsPluralBucket(%q) = true, want false", b)
		}
	}
	if IsPluralBucket("s") {
		t.Errorf("IsPluralBucket(\"s\") = true, want false (single char)")
	}
}

func TestListLengthKey(t *testing.T) {
	if got := ListLengthKey("servers"); got != "servers/:len" {
		t.Errorf("ListLengthKey(\"servers\") = %q, want servers/:len", got)
	}
}

func TestListItemKey(t *testing.T) {
	if got := ListItemKey("servers", 0); got != "servers/:0" {
		t.Errorf("ListItemKey(\"servers\", 0) = %q, want servers/:0", got)
	}
	if got := ListItemKey("servers/personal", 0); got != "servers/personal/:0" {
		t.Errorf("ListItemKey(\"servers/personal\", 0) = %q, want servers/personal/:0", got)
	}
}
