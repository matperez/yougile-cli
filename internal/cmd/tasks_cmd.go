package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewTasksListCmd returns the "tasks list" command.
func NewTasksListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title, columnID string

	c := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.TaskControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}
			if columnID != "" {
				params.ColumnId = strPtr(columnID)
			}

			resp, err := api.TaskControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list tasks: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list tasks: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list tasks: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title", "ColumnId"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, t := range resp.JSON200.Content {
				colID := ""
				if t.ColumnId != nil {
					colID = *t.ColumnId
				}
				rows = append(rows, []string{t.Id, t.Title, colID})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	c.Flags().StringVar(&columnID, "column-id", "", "filter by column ID")
	return c
}

// NewTaskGetCmd returns the "tasks get" command.
func NewTaskGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get task by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.TaskControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get task: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get task: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get task: empty response")
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

// NewTasksCmd returns the "tasks" parent command.
func NewTasksCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "tasks",
		Short: "Manage tasks",
	}
	c.AddCommand(NewTasksListCmd(resolvePath, outputJSON))
	c.AddCommand(NewTaskGetCmd(resolvePath, outputJSON))
	return c
}
