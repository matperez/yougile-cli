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
	return c
}
