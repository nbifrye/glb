package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HostConfig struct {
	Token       string `yaml:"token"`
	APIProtocol string `yaml:"api_protocol,omitempty"`
}

type configData struct {
	Hosts map[string]*HostConfig `yaml:"hosts"`
}

type Config struct {
	data     configData
	filePath string
}

func ConfigDir() string {
	if dir := os.Getenv("GLB_CONFIG_DIR"); dir != "" {
		return dir
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "glb")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "glb")
}

func configFilePath() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

func Load() (*Config, error) {
	c := &Config{
		data:     configData{Hosts: make(map[string]*HostConfig)},
		filePath: configFilePath(),
	}

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &c.data); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if c.data.Hosts == nil {
		c.data.Hosts = make(map[string]*HostConfig)
	}
	return c, nil
}

func (c *Config) Hosts() []string {
	hosts := make([]string, 0, len(c.data.Hosts))
	for h := range c.data.Hosts {
		hosts = append(hosts, h)
	}
	return hosts
}

func (c *Config) Token(hostname string) string {
	if h, ok := c.data.Hosts[hostname]; ok {
		return h.Token
	}
	return ""
}

func (c *Config) APIProtocol(hostname string) string {
	if h, ok := c.data.Hosts[hostname]; ok && h.APIProtocol != "" {
		return h.APIProtocol
	}
	return "https"
}

func (c *Config) SetHost(hostname, token, apiProtocol string) {
	c.data.Hosts[hostname] = &HostConfig{
		Token:       token,
		APIProtocol: apiProtocol,
	}
}

func (c *Config) Write() error {
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(&c.data)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(c.filePath, data, 0o600)
}
