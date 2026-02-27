package cmd

import (
	"context"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/auth"
	"github.com/angolovin/yougile-cli/internal/config"
	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
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

// NewAuthCompaniesCmd returns the "auth companies" command (list companies by email/password).
func NewAuthCompaniesCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var email, password string
	c := &cobra.Command{
		Use:   "companies",
		Short: "List companies (requires email and password)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" || password == "" {
				return fmt.Errorf("email and password are required (--email, --password)")
			}
			path, err := resolvePath()
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			api, err := client.NewClientWithResponses(cfg.BaseURL)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			resp, err := api.GetCompaniesWithResponse(context.Background(), nil, client.GetCompaniesJSONRequestBody{
				Login:    email,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("get companies: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get companies: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get companies: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Name", "Admin"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, co := range resp.JSON200.Content {
				admin := "no"
				if co.IsAdmin {
					admin = "yes"
				}
				rows = append(rows, []string{co.Id, co.Name, admin})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().StringVar(&email, "email", "", "account email")
	c.Flags().StringVar(&password, "password", "", "account password")
	_ = c.MarkFlagRequired("email")
	_ = c.MarkFlagRequired("password")
	return c
}

// NewAuthKeysListCmd returns the "auth keys list" command.
func NewAuthKeysListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var email, password, companyID string
	c := &cobra.Command{
		Use:   "list",
		Short: "List API keys (requires email and password)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" || password == "" {
				return fmt.Errorf("email and password are required (--email, --password)")
			}
			path, err := resolvePath()
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			api, err := client.NewClientWithResponses(cfg.BaseURL)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			body := client.AuthKeyControllerSearchJSONRequestBody{Login: email, Password: password}
			if companyID != "" {
				body.CompanyId = &companyID
			}
			resp, err := api.AuthKeyControllerSearchWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("list keys: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list keys: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list keys: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			keys := resp.JSON200
			if keys == nil {
				keys = &[]client.AuthKeyWithDetailsDto{}
			}
			headers := []string{"Key", "CompanyId", "Deleted"}
			rows := make([][]string, 0, len(*keys))
			for _, k := range *keys {
				del := "no"
				if k.Deleted {
					del = "yes"
				}
				rows = append(rows, []string{k.Key, k.CompanyId, del})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().StringVar(&email, "email", "", "account email")
	c.Flags().StringVar(&password, "password", "", "account password")
	c.Flags().StringVar(&companyID, "company-id", "", "filter by company ID")
	_ = c.MarkFlagRequired("email")
	_ = c.MarkFlagRequired("password")
	return c
}

// NewAuthKeysCreateCmd returns the "auth keys create" command.
func NewAuthKeysCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var email, password, companyID string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create an API key (requires email, password, company-id)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" || password == "" || companyID == "" {
				return fmt.Errorf("email, password and company-id are required")
			}
			path, err := resolvePath()
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			api, err := client.NewClientWithResponses(cfg.BaseURL)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			body := client.AuthKeyControllerCreateJSONRequestBody{
				Login:     email,
				Password:  password,
				CompanyId: companyID,
			}
			resp, err := api.AuthKeyControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create key: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create key: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "API key created: %s\n", resp.JSON201.Key)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&email, "email", "", "account email")
	c.Flags().StringVar(&password, "password", "", "account password")
	c.Flags().StringVar(&companyID, "company-id", "", "company ID")
	_ = c.MarkFlagRequired("email")
	_ = c.MarkFlagRequired("password")
	_ = c.MarkFlagRequired("company-id")
	return c
}

// NewAuthKeysDeleteCmd returns the "auth keys delete" command.
func NewAuthKeysDeleteCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [key]",
		Short: "Delete an API key by key value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolvePath()
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			api, err := client.NewClientWithResponses(cfg.BaseURL)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			key := args[0]
			resp, err := api.AuthKeyControllerDeleteWithResponse(context.Background(), key)
			if err != nil {
				return fmt.Errorf("delete key: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("delete key: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			_, err = fmt.Fprintf(out, "API key deleted\n")
			return err
		},
	}
}

// NewAuthKeysCmd returns the "auth keys" parent command.
func NewAuthKeysCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "keys",
		Short: "Manage API keys",
	}
	c.AddCommand(NewAuthKeysListCmd(resolvePath, outputJSON))
	c.AddCommand(NewAuthKeysCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewAuthKeysDeleteCmd(resolvePath, outputJSON))
	return c
}

// NewAuthCmd returns the "auth" parent command with login, companies, keys.
func NewAuthCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: "Authentication and API keys",
	}
	c.AddCommand(NewAuthLoginCmd(resolvePath))
	c.AddCommand(NewAuthCompaniesCmd(resolvePath, outputJSON))
	c.AddCommand(NewAuthKeysCmd(resolvePath, outputJSON))
	return c
}
