package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	t.Run("GLB_CONFIG_DIR takes precedence", func(t *testing.T) {
		t.Setenv("GLB_CONFIG_DIR", "/custom/dir")
		t.Setenv("XDG_CONFIG_HOME", "/xdg")
		if got := ConfigDir(); got != "/custom/dir" {
			t.Errorf("ConfigDir() = %q, want %q", got, "/custom/dir")
		}
	})

	t.Run("XDG_CONFIG_HOME fallback", func(t *testing.T) {
		t.Setenv("GLB_CONFIG_DIR", "")
		t.Setenv("XDG_CONFIG_HOME", "/xdg")
		want := filepath.Join("/xdg", "glb")
		if got := ConfigDir(); got != want {
			t.Errorf("ConfigDir() = %q, want %q", got, want)
		}
	})
}

func TestLoadWrite(t *testing.T) {
	// Use a subdirectory inside TempDir so MkdirAll creates it with 0o700
	dir := filepath.Join(t.TempDir(), "glb-config")
	t.Setenv("GLB_CONFIG_DIR", dir)

	// Load should succeed even if file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Hosts()) != 0 {
		t.Errorf("empty config should have 0 hosts, got %d", len(cfg.Hosts()))
	}

	// Set a host and write
	cfg.SetHost("gitlab.example.com", "test-token", "https")
	if err := cfg.Write(); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify directory permissions
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("config dir permissions = %o, want %o", perm, 0o700)
	}

	// Re-load and verify
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load() after write error = %v", err)
	}
	if got := cfg2.Token("gitlab.example.com"); got != "test-token" {
		t.Errorf("Token() = %q, want %q", got, "test-token")
	}
	if got := cfg2.APIProtocol("gitlab.example.com"); got != "https" {
		t.Errorf("APIProtocol() = %q, want %q", got, "https")
	}
}

func TestAPIProtocolDefault(t *testing.T) {
	cfg := &Config{data: configData{Hosts: make(map[string]*HostConfig)}}
	if got := cfg.APIProtocol("unknown.host"); got != "https" {
		t.Errorf("APIProtocol(unknown) = %q, want %q", got, "https")
	}
}

func TestTokenNotFound(t *testing.T) {
	cfg := &Config{data: configData{Hosts: make(map[string]*HostConfig)}}
	if got := cfg.Token("unknown.host"); got != "" {
		t.Errorf("Token(unknown) = %q, want empty", got)
	}
}
