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
