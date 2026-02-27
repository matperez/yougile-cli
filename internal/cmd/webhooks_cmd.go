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

// NewWebhooksCreateCmd returns the "webhooks create" command.
func NewWebhooksCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var event, url string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == "" || url == "" {
				return fmt.Errorf("event and url are required (--event, --url)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.WebhookControllerCreateJSONRequestBody{
				Event:   event,
				Url:     url,
				Filters: []client.WebhookFilters{},
			}
			resp, err := api.WebhookControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create webhook: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create webhook: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Webhook created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&event, "event", "", "event pattern (e.g. task-*, .*)")
	c.Flags().StringVar(&url, "url", "", "webhook URL")
	_ = c.MarkFlagRequired("event")
	_ = c.MarkFlagRequired("url")
	return c
}

// NewWebhooksCmd returns the "webhooks" parent command.
func NewWebhooksCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage webhooks",
	}
	c.AddCommand(NewWebhooksListCmd(resolvePath, outputJSON))
	c.AddCommand(NewWebhooksCreateCmd(resolvePath, outputJSON))
	return c
}
