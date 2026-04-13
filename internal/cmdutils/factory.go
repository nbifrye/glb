package cmdutils

import (
	"fmt"

	"github.com/nbifrye/glb/internal/api"
	"github.com/nbifrye/glb/internal/auth"
	"github.com/nbifrye/glb/internal/config"
	"github.com/nbifrye/glb/internal/iostreams"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Factory struct {
	Config func() (*config.Config, error)
	IO     *iostreams.IOStreams
}

func NewFactory() *Factory {
	var cachedConfig *config.Config
	return &Factory{
		Config: func() (*config.Config, error) {
			if cachedConfig != nil {
				return cachedConfig, nil
			}
			cfg, err := config.Load()
			if err != nil {
				return nil, err
			}
			cachedConfig = cfg
			return cfg, nil
		},
		IO: iostreams.System(),
	}
}

func (f *Factory) GitLabClient() (*gitlab.Client, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}

	hostname, token := auth.DefaultHostWithToken(cfg)
	if token == "" {
		return nil, fmt.Errorf("authentication required: run 'glb auth login' or set GITLAB_TOKEN")
	}

	protocol := cfg.APIProtocol(hostname)
	return api.NewGitLabClient(token, hostname, protocol)
}

func (f *Factory) Hostname() (string, error) {
	cfg, err := f.Config()
	if err != nil {
		return "", err
	}
	hostname, _ := auth.DefaultHostWithToken(cfg)
	return hostname, nil
}
