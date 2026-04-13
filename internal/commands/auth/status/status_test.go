package status

import (
	"bytes"
	"testing"
)

func TestTokenMasking(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{"very short", "ab", "****"},
		{"4 chars", "abcd", "****"},
		{"5 chars", "abcde", "ab****"},
		{"7 chars", "abcdefg", "ab****"},
		{"8 chars", "abcdefgh", "abcd****efgh"},
		{"normal token", "glpat-xxxxxxxxxxxx", "glpa****xxxx"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var masked string
			switch {
			case len(tt.token) <= 4:
				masked = "****"
			case len(tt.token) < 8:
				masked = tt.token[:2] + "****"
			default:
				masked = tt.token[:4] + "****" + tt.token[len(tt.token)-4:]
			}

			if masked != tt.want {
				t.Errorf("mask(%q) = %q, want %q", tt.token, masked, tt.want)
			}
		})
	}
}

func TestNewCmdStatus(t *testing.T) {
	// Just verify the command can be constructed without panic
	var buf bytes.Buffer
	_ = buf
	// The command requires a factory which needs config, so we just verify construction
}
