package cmd

import (
	"context"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/auth"
	"github.com/angolovin/yougile-cli/internal/config"
	"github.com/spf13/cobra"
)

// NewAuthLoginCmd returns the "auth login" command.
func NewAuthLoginCmd(resolvePath func() (string, error)) *cobra.Command {
	var email, password string

	c := &cobra.Command{
		Use:   "login",
		Short: "Log in with email and password, save API key to config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" || password == "" {
				return fmt.Errorf("email and password are required (use --email and --password)")
			}

			path, err := resolvePath()
			if err != nil {
				return fmt.Errorf("resolve config path: %w", err)
			}

			key, err := auth.Login(context.Background(), config.DefaultBaseURL(), email, password)
			if err != nil {
				return fmt.Errorf("login: %w", err)
			}

			cfg := &config.Config{
				BaseURL: config.DefaultBaseURL(),
				APIKey:  key,
			}
			if err := config.Save(path, cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "API key saved to %s\n", path)
			return nil
		},
	}
	c.Flags().StringVar(&email, "email", "", "account email")
	c.Flags().StringVar(&password, "password", "", "account password")
	_ = c.MarkFlagRequired("email")
	_ = c.MarkFlagRequired("password")
	return c
}

// NewAuthCmd returns the "auth" parent command with login subcommand.
func NewAuthCmd(resolvePath func() (string, error)) *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: "Authentication and API keys",
	}
	c.AddCommand(NewAuthLoginCmd(resolvePath))
	return c
}
