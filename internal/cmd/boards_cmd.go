package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewBoardsListCmd returns the "boards list" command.
func NewBoardsListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title, projectID string

	c := &cobra.Command{
		Use:   "list",
		Short: "List boards",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}

			params := &client.BoardControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}
			if projectID != "" {
				params.ProjectId = strPtr(projectID)
			}

			resp, err := api.BoardControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list boards: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list boards: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list boards: empty response")
			}

			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title", "ProjectId"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, b := range resp.JSON200.Content {
				rows = append(rows, []string{b.Id, b.Title, b.ProjectId})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	c.Flags().StringVar(&projectID, "project-id", "", "filter by project ID")
	return c
}

// NewBoardsCreateCmd returns the "boards create" command.
func NewBoardsCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, projectID string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a board",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" || projectID == "" {
				return fmt.Errorf("title and project-id are required")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.BoardControllerCreateJSONRequestBody{Title: title, ProjectId: projectID}
			resp, err := api.BoardControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create board: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create board: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Board created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "board title")
	c.Flags().StringVar(&projectID, "project-id", "", "project ID")
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("project-id")
	return c
}

// NewBoardsUpdateCmd returns the "boards update" command.
func NewBoardsUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a board",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.BoardControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("title") {
				body.Title = &title
			}
			resp, err := api.BoardControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update board: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update board: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Board updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&title, "title", "", "board title")
	return c
}

// NewBoardGetCmd returns the "boards get" command.
func NewBoardGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get board by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]

			resp, err := api.BoardControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get board: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get board: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get board: empty response")
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

// NewBoardsCmd returns the "boards" parent command.
func NewBoardsCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "boards",
		Short: "Manage boards",
	}
	c.AddCommand(NewBoardsListCmd(resolvePath, outputJSON))
	c.AddCommand(NewBoardGetCmd(resolvePath, outputJSON))
	c.AddCommand(NewBoardsCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewBoardsUpdateCmd(resolvePath, outputJSON))
	return c
}
