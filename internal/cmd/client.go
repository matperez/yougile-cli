package cmd

import (
	"context"
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
		return nil, nil
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
