package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewWebhooksListCmd returns the "webhooks list" command.
func NewWebhooksListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var includeDeleted bool
	c := &cobra.Command{
		Use:   "list",
		Short: "List webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.WebhookControllerSearchParams{}
			if includeDeleted {
				params.IncludeDeleted = &includeDeleted
			}
			resp, err := api.WebhookControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list webhooks: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list webhooks: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			if resp.JSON200 == nil {
				_, err = fmt.Fprintln(out, "{}")
				return err
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(resp.JSON200)
		},
	}
	c.Flags().BoolVar(&includeDeleted, "include-deleted", false, "include deleted webhooks")
	return c
}

// NewWebhooksCmd returns the "webhooks" parent command.
func NewWebhooksCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage webhooks",
	}
	c.AddCommand(NewWebhooksListCmd(resolvePath, outputJSON))
	return c
}
