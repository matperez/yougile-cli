package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/angolovin/yougile-cli/internal/config"
	"github.com/angolovin/yougile-cli/pkg/client"
)

// NewAPIClient returns a YouGile API client with Bearer auth from cfg.
// baseURL is normalized (no trailing slash).
func NewAPIClient(cfg *config.Config) (*client.ClientWithResponses, error) {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key not set in config")
	}
	apiKey := cfg.APIKey
	opts := []client.ClientOption{
		client.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apiKey)
			return nil
		}),
	}
	return client.NewClientWithResponses(baseURL, opts...)
}

// loadConfigAndClient loads config from path and creates API client.
// Returns error if config missing, invalid, or api_key empty.
func loadConfigAndClient(resolvePath func() (string, error)) (*config.Config, *client.ClientWithResponses, error) {
	path, err := resolvePath()
	if err != nil {
		return nil, nil, fmt.Errorf("resolve config path: %w", err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("load config: %w", err)
	}
	api, err := NewAPIClient(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("create API client: %w", err)
	}
	return cfg, api, nil
}
