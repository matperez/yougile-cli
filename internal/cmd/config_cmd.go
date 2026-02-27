package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/config"
	"github.com/spf13/cobra"
)

const apiKeyMask = "***"

// NewConfigPathCmd returns the "config path" command.
func NewConfigPathCmd(resolvePath func() (string, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print path to config file",
		RunE: func(c *cobra.Command, args []string) error {
			path, err := resolvePath()
			if err != nil {
				return fmt.Errorf("resolve config path: %w", err)
			}
			_, writeErr := fmt.Fprintln(c.OutOrStdout(), path)
			if writeErr != nil {
				return writeErr
			}
			return nil
		},
	}
}

// NewConfigShowCmd returns the "config show" command.
func NewConfigShowCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current config (api_key masked in human output)",
		RunE: func(c *cobra.Command, args []string) error {
			path, err := resolvePath()
			if err != nil {
				return fmt.Errorf("resolve config path: %w", err)
			}

			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			out := c.OutOrStdout()
			if outputJSON() {
				return json.NewEncoder(out).Encode(cfg)
			}

			// Human output: mask api_key
			display := struct {
				BaseURL string `json:"base_url"`
				APIKey  string `json:"api_key"`
			}{
				BaseURL: cfg.BaseURL,
				APIKey:  apiKeyMask,
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(display)
		},
	}
}

// NewConfigCmd returns the "config" parent command with path and show subcommands.
func NewConfigCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Manage config file",
	}
	c.AddCommand(NewConfigPathCmd(resolvePath))
	c.AddCommand(NewConfigShowCmd(resolvePath, outputJSON))
	return c
}
