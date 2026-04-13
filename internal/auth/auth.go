package auth

import (
	"os"

	"github.com/nbifrye/glb/internal/config"
	"github.com/nbifrye/glb/internal/glinstance"
)

func TokenForHost(cfg *config.Config, hostname string) string {
	if token := os.Getenv("GITLAB_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GLB_TOKEN"); token != "" {
		return token
	}
	return cfg.Token(hostname)
}

func DefaultHostWithToken(cfg *config.Config) (hostname, token string) {
	if token := os.Getenv("GITLAB_TOKEN"); token != "" {
		return hostFromEnvOrConfig(cfg), token
	}
	if token := os.Getenv("GLB_TOKEN"); token != "" {
		return hostFromEnvOrConfig(cfg), token
	}

	for _, h := range cfg.Hosts() {
		if t := cfg.Token(h); t != "" {
			return h, t
		}
	}
	return glinstance.Default(), ""
}

func hostFromEnvOrConfig(cfg *config.Config) string {
	if h := os.Getenv("GITLAB_HOST"); h != "" {
		return glinstance.NormalizeHostname(h)
	}
	if h := os.Getenv("GLB_HOST"); h != "" {
		return glinstance.NormalizeHostname(h)
	}
	hosts := cfg.Hosts()
	if len(hosts) == 1 {
		return hosts[0]
	}
	return glinstance.Default()
}
