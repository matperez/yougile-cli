package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultBaseURL = "https://ru.yougile.com"

// Config holds YouGile CLI configuration.
type Config struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
}

// Load reads and parses the config file at path.
// Returns error if the file doesn't exist or is invalid.
// If base_url is empty, it is set to defaultBaseURL.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	return &cfg, nil
}
