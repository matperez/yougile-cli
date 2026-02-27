package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewProjectsListCmd returns the "projects list" command.
func NewProjectsListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title string

	c := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.ProjectControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}

			resp, err := api.ProjectControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list projects: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list projects: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list projects: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, p := range resp.JSON200.Content {
				rows = append(rows, []string{p.Id, p.Title})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	return c
}

// NewProjectsCreateCmd returns the "projects create" command.
func NewProjectsCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("title is required (--title)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.ProjectControllerCreateJSONRequestBody{Title: title}
			resp, err := api.ProjectControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create project: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create project: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Project created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "project title")
	_ = c.MarkFlagRequired("title")
	return c
}

// NewProjectsUpdateCmd returns the "projects update" command.
func NewProjectsUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.ProjectControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("title") {
				body.Title = &title
			}
			resp, err := api.ProjectControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update project: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update project: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Project updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&title, "title", "", "project title")
	return c
}

// NewProjectGetCmd returns the "projects get" command.
func NewProjectGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get project by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.ProjectControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get project: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get project: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get project: empty response")
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

// NewProjectRolesListCmd returns the "projects roles list" command.
func NewProjectRolesListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var projectID string
	var limit, offset int
	c := &cobra.Command{
		Use:   "list",
		Short: "List project roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("project-id is required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.ProjectRolesControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			resp, err := api.ProjectRolesControllerSearchWithResponse(context.Background(), projectID, params)
			if err != nil {
				return fmt.Errorf("list roles: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list roles: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list roles: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Name"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, r := range resp.JSON200.Content {
				rows = append(rows, []string{r.Id, r.Name})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	c.Flags().IntVar(&limit, "limit", 50, "max items")
	c.Flags().IntVar(&offset, "offset", 0, "offset")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewProjectRolesGetCmd returns the "projects roles get" command.
func NewProjectRolesGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var projectID string
	c := &cobra.Command{
		Use:   "get [role-id]",
		Short: "Get project role by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("project-id is required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			roleID := args[0]
			resp, err := api.ProjectRolesControllerGetWithResponse(context.Background(), projectID, roleID)
			if err != nil {
				return fmt.Errorf("get role: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get role: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get role: empty response")
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
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewProjectRolesCreateCmd returns the "projects roles create" command.
func NewProjectRolesCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var projectID, name, description string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a project role",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" || name == "" {
				return fmt.Errorf("project-id and name are required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.ProjectRolesControllerCreateJSONRequestBody{
				Name:        name,
				Permissions: client.ProjectPermissionsDto{}, // minimal permissions
			}
			if description != "" {
				body.Description = &description
			}
			resp, err := api.ProjectRolesControllerCreateWithResponse(context.Background(), projectID, body)
			if err != nil {
				return fmt.Errorf("create role: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create role: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Role created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	c.Flags().StringVar(&name, "name", "", "role name")
	c.Flags().StringVar(&description, "description", "", "role description")
	_ = c.MarkFlagRequired("project-id")
	_ = c.MarkFlagRequired("name")
	return c
}

// NewProjectRolesUpdateCmd returns the "projects roles update" command.
func NewProjectRolesUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var projectID, name, description string
	c := &cobra.Command{
		Use:   "update [role-id]",
		Short: "Update a project role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("project-id is required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			roleID := args[0]
			body := client.ProjectRolesControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("name") {
				body.Name = &name
			}
			if cmd.Flags().Changed("description") {
				body.Description = &description
			}
			resp, err := api.ProjectRolesControllerUpdateWithResponse(context.Background(), projectID, roleID, body)
			if err != nil {
				return fmt.Errorf("update role: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update role: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Role updated: id=%s\n", roleID)
			return err
		},
	}
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	c.Flags().StringVar(&name, "name", "", "role name")
	c.Flags().StringVar(&description, "description", "", "role description")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewProjectRolesDeleteCmd returns the "projects roles delete" command.
func NewProjectRolesDeleteCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var projectID string
	c := &cobra.Command{
		Use:   "delete [role-id]",
		Short: "Delete a project role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("project-id is required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			roleID := args[0]
			resp, err := api.ProjectRolesControllerDeleteWithResponse(context.Background(), projectID, roleID)
			if err != nil {
				return fmt.Errorf("delete role: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("delete role: HTTP %s", resp.HTTPResponse.Status)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Role deleted: id=%s\n", roleID)
			return err
		},
	}
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewProjectsCmd returns the "projects" parent command.
func NewProjectsCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "projects",
		Short: "Manage projects",
	}
	c.AddCommand(NewProjectsListCmd(resolvePath, outputJSON))
	c.AddCommand(NewProjectGetCmd(resolvePath, outputJSON))
	c.AddCommand(NewProjectsCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewProjectsUpdateCmd(resolvePath, outputJSON))
	rolesCmd := &cobra.Command{Use: "roles", Short: "Project roles"}
	rolesCmd.AddCommand(NewProjectRolesListCmd(resolvePath, outputJSON))
	rolesCmd.AddCommand(NewProjectRolesGetCmd(resolvePath, outputJSON))
	rolesCmd.AddCommand(NewProjectRolesCreateCmd(resolvePath, outputJSON))
	rolesCmd.AddCommand(NewProjectRolesUpdateCmd(resolvePath, outputJSON))
	rolesCmd.AddCommand(NewProjectRolesDeleteCmd(resolvePath, outputJSON))
	c.AddCommand(rolesCmd)
	return c
}
