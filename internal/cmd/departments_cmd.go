package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewDepartmentsListCmd returns the "departments list" command.
func NewDepartmentsListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title, parentID string

	c := &cobra.Command{
		Use:   "list",
		Short: "List departments",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.DepartmentControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}
			if parentID != "" {
				params.ParentId = strPtr(parentID)
			}

			resp, err := api.DepartmentControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list departments: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list departments: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list departments: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title", "ParentId"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, d := range resp.JSON200.Content {
				pID := ""
				if d.ParentId != nil {
					pID = *d.ParentId
				}
				rows = append(rows, []string{d.Id, d.Title, pID})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	c.Flags().StringVar(&parentID, "parent-id", "", "filter by parent department ID")
	return c
}

// NewDepartmentsCreateCmd returns the "departments create" command.
func NewDepartmentsCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, parentID string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a department",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("title is required (--title)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.DepartmentControllerCreateJSONRequestBody{Title: title}
			if parentID != "" {
				body.ParentId = &parentID
			}
			resp, err := api.DepartmentControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create department: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create department: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Department created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "department title")
	c.Flags().StringVar(&parentID, "parent-id", "", "parent department ID (empty for top-level)")
	_ = c.MarkFlagRequired("title")
	return c
}

// NewDepartmentsUpdateCmd returns the "departments update" command.
func NewDepartmentsUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a department",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.DepartmentControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("title") {
				body.Title = &title
			}
			resp, err := api.DepartmentControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update department: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update department: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Department updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&title, "title", "", "department title")
	return c
}

// NewDepartmentGetCmd returns the "departments get" command.
func NewDepartmentGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get department by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.DepartmentControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get department: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get department: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get department: empty response")
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

// NewDepartmentsCmd returns the "departments" parent command.
func NewDepartmentsCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "departments",
		Short: "Manage departments",
	}
	c.AddCommand(NewDepartmentsListCmd(resolvePath, outputJSON))
	c.AddCommand(NewDepartmentGetCmd(resolvePath, outputJSON))
	c.AddCommand(NewDepartmentsCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewDepartmentsUpdateCmd(resolvePath, outputJSON))
	return c
}
