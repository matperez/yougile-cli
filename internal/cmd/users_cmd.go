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
func boolPtr(b bool) *bool           { return &b }

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

// NewUsersCreateCmd returns the "users create" command.
func NewUsersCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var email string
	var isAdmin bool
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" {
				return fmt.Errorf("email is required (--email)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.UserControllerCreateJSONRequestBody{
				Email:   email,
				IsAdmin: &isAdmin,
			}
			resp, err := api.UserControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create user: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create user: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "User created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&email, "email", "", "user email")
	c.Flags().BoolVar(&isAdmin, "admin", false, "grant admin rights")
	_ = c.MarkFlagRequired("email")
	return c
}

// NewUsersUpdateCmd returns the "users update" command.
func NewUsersUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var isAdmin bool
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.UserControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("admin") {
				body.IsAdmin = &isAdmin
			}
			resp, err := api.UserControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update user: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update user: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			if resp.JSON200 != nil {
				_, err = fmt.Fprintf(out, "User updated: id=%s\n", id)
				return err
			}
			return nil
		},
	}
	c.Flags().BoolVar(&isAdmin, "admin", false, "set admin rights")
	return c
}

// NewUsersDeleteCmd returns the "users delete" command.
func NewUsersDeleteCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			resp, err := api.UserControllerDeleteWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("delete user: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("delete user: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "User deleted: id=%s\n", id)
			return err
		},
	}
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
	c.AddCommand(NewUsersCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewUsersUpdateCmd(resolvePath, outputJSON))
	c.AddCommand(NewUsersDeleteCmd(resolvePath, outputJSON))
	return c
}
