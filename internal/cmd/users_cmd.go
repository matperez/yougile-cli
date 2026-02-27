package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

func float32Ptr(v float32) *float32 { return &v }
func strPtr(s string) *string       { return &s }

// NewUsersListCmd returns the "users list" command.
func NewUsersListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var email, projectID string

	c := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.UserControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if email != "" {
				params.Email = strPtr(email)
			}
			if projectID != "" {
				params.ProjectId = strPtr(projectID)
			}

			resp, err := api.UserControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list users: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list users: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list users: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Email", "Admin"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, u := range resp.JSON200.Content {
				admin := "no"
				if u.IsAdmin != nil && *u.IsAdmin {
					admin = "yes"
				}
				rows = append(rows, []string{u.Id, u.Email, admin})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&email, "email", "", "filter by email")
	c.Flags().StringVar(&projectID, "project-id", "", "filter by project ID")
	return c
}

// NewUserGetCmd returns the "users get" command.
func NewUserGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get user by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.UserControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get user: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get user: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get user: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(resp.JSON200)
		},
	}
}

// NewUsersCmd returns the "users" parent command.
func NewUsersCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}
	c.AddCommand(NewUsersListCmd(resolvePath, outputJSON))
	c.AddCommand(NewUserGetCmd(resolvePath, outputJSON))
	return c
}
