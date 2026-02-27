package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewColumnsListCmd returns the "columns list" command.
func NewColumnsListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title, boardID string

	c := &cobra.Command{
		Use:   "list",
		Short: "List columns",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.ColumnControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}
			if boardID != "" {
				params.BoardId = strPtr(boardID)
			}

			resp, err := api.ColumnControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list columns: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list columns: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list columns: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title", "BoardId"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, col := range resp.JSON200.Content {
				rows = append(rows, []string{col.Id, col.Title, col.BoardId})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	c.Flags().StringVar(&boardID, "board-id", "", "filter by board ID")
	return c
}

// NewColumnsCreateCmd returns the "columns create" command.
func NewColumnsCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, boardID string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a column",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" || boardID == "" {
				return fmt.Errorf("title and board-id are required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.ColumnControllerCreateJSONRequestBody{Title: title, BoardId: boardID}
			resp, err := api.ColumnControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create column: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create column: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Column created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "column title")
	c.Flags().StringVar(&boardID, "board-id", "", "board ID")
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("board-id")
	return c
}

// NewColumnsUpdateCmd returns the "columns update" command.
func NewColumnsUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a column",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.ColumnControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("title") {
				body.Title = &title
			}
			resp, err := api.ColumnControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update column: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update column: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Column updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&title, "title", "", "column title")
	return c
}

// NewColumnGetCmd returns the "columns get" command.
func NewColumnGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get column by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.ColumnControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get column: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get column: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get column: empty response")
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

// NewColumnsCmd returns the "columns" parent command.
func NewColumnsCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "columns",
		Short: "Manage columns",
	}
	c.AddCommand(NewColumnsListCmd(resolvePath, outputJSON))
	c.AddCommand(NewColumnGetCmd(resolvePath, outputJSON))
	c.AddCommand(NewColumnsCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewColumnsUpdateCmd(resolvePath, outputJSON))
	return c
}
