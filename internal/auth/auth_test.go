package auth

import (
	"testing"

	"github.com/nbifrye/glb/internal/config"
)

func newTestConfig(hosts map[string]string) *config.Config {
	cfg := config.NewEmpty()
	for host, token := range hosts {
		cfg.SetHost(host, token, "https")
	}
	return cfg
}

func TestTokenForHost(t *testing.T) {
	t.Run("GITLAB_TOKEN takes precedence", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "env-token")
		t.Setenv("GLB_TOKEN", "glb-token")
		cfg := newTestConfig(map[string]string{"gitlab.com": "config-token"})
		if got := TokenForHost(cfg, "gitlab.com"); got != "env-token" {
			t.Errorf("TokenForHost() = %q, want %q", got, "env-token")
		}
	})

	t.Run("GLB_TOKEN fallback", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "")
		t.Setenv("GLB_TOKEN", "glb-token")
		cfg := newTestConfig(map[string]string{"gitlab.com": "config-token"})
		if got := TokenForHost(cfg, "gitlab.com"); got != "glb-token" {
			t.Errorf("TokenForHost() = %q, want %q", got, "glb-token")
		}
	})

	t.Run("config fallback", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "")
		t.Setenv("GLB_TOKEN", "")
		cfg := newTestConfig(map[string]string{"gitlab.com": "config-token"})
		if got := TokenForHost(cfg, "gitlab.com"); got != "config-token" {
			t.Errorf("TokenForHost() = %q, want %q", got, "config-token")
		}
	})
}

func TestDefaultHostWithToken(t *testing.T) {
	t.Run("env var with GITLAB_HOST", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "env-token")
		t.Setenv("GITLAB_HOST", "custom.gitlab.com")
		t.Setenv("GLB_HOST", "")
		cfg := newTestConfig(nil)
		host, token := DefaultHostWithToken(cfg)
		if host != "custom.gitlab.com" {
			t.Errorf("host = %q, want %q", host, "custom.gitlab.com")
		}
		if token != "env-token" {
			t.Errorf("token = %q, want %q", token, "env-token")
		}
	})

	t.Run("env var with single config host", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "env-token")
		t.Setenv("GITLAB_HOST", "")
		t.Setenv("GLB_HOST", "")
		cfg := newTestConfig(map[string]string{"my.gitlab.io": "ignored"})
		host, token := DefaultHostWithToken(cfg)
		if host != "my.gitlab.io" {
			t.Errorf("host = %q, want %q", host, "my.gitlab.io")
		}
		if token != "env-token" {
			t.Errorf("token = %q, want %q", token, "env-token")
		}
	})

	t.Run("no env var uses config host", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "")
		t.Setenv("GLB_TOKEN", "")
		t.Setenv("GITLAB_HOST", "")
		t.Setenv("GLB_HOST", "")
		cfg := newTestConfig(map[string]string{"custom.host": "my-token"})
		host, token := DefaultHostWithToken(cfg)
		if host != "custom.host" {
			t.Errorf("host = %q, want %q", host, "custom.host")
		}
		if token != "my-token" {
			t.Errorf("token = %q, want %q", token, "my-token")
		}
	})

	t.Run("no config returns default", func(t *testing.T) {
		t.Setenv("GITLAB_TOKEN", "")
		t.Setenv("GLB_TOKEN", "")
		t.Setenv("GITLAB_HOST", "")
		t.Setenv("GLB_HOST", "")
		cfg := newTestConfig(nil)
		host, token := DefaultHostWithToken(cfg)
		if host != "gitlab.com" {
			t.Errorf("host = %q, want %q", host, "gitlab.com")
		}
		if token != "" {
			t.Errorf("token = %q, want empty", token)
		}
	})
}
