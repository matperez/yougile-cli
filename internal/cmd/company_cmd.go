package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/config"
	"github.com/spf13/cobra"
)

// NewCompanyGetCmd returns the "company get" command.
func NewCompanyGetCmd(
	resolvePath func() (string, error),
	outputJSON func() bool,
) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current company details",
		RunE: func(c *cobra.Command, args []string) error {
			path, err := resolvePath()
			if err != nil {
				return fmt.Errorf("resolve config path: %w", err)
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if cfg.APIKey == "" {
				return fmt.Errorf("api_key not set in config; run 'yougile auth login' first")
			}

			api, err := NewAPIClient(cfg)
			if err != nil {
				return fmt.Errorf("create API client: %w", err)
			}

			resp, err := api.CompanyControllerGetWithResponse(context.Background())
			if err != nil {
				return fmt.Errorf("get company: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get company: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get company: empty response")
			}

			out := c.OutOrStdout()
			if outputJSON() {
				return json.NewEncoder(out).Encode(resp.JSON200)
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(resp.JSON200)
		},
	}
}

// NewCompanyCmd returns the "company" parent command with get subcommand.
func NewCompanyCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	cc := &cobra.Command{
		Use:   "company",
		Short: "Company details",
	}
	cc.AddCommand(NewCompanyGetCmd(resolvePath, outputJSON))
	return cc
}
