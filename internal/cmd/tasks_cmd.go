package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

// NewTasksCreateCmd returns the "tasks create" command.
func NewTasksCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, columnID string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("title is required (--title)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.TaskControllerCreateJSONRequestBody{Title: title}
			if columnID != "" {
				body.ColumnId = &columnID
			}
			resp, err := api.TaskControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create task: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create task: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Task created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "task title")
	c.Flags().StringVar(&columnID, "column-id", "", "column ID (optional)")
	_ = c.MarkFlagRequired("title")
	return c
}

// NewTasksUpdateCmd returns the "tasks update" command.
func NewTasksUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var title, columnID string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.TaskControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("title") {
				body.Title = &title
			}
			if cmd.Flags().Changed("column-id") {
				body.ColumnId = &columnID
			}
			resp, err := api.TaskControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update task: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update task: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Task updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&title, "title", "", "task title")
	c.Flags().StringVar(&columnID, "column-id", "", "column ID (move task to another column)")
	return c
}

// NewTasksChatSubscribersGetCmd returns the "tasks chat-subscribers get" command.
func NewTasksChatSubscribersGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [task-id]",
		Short: "Get task chat subscribers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			resp, err := api.TaskControllerGetChatSubscribersWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get chat subscribers: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get chat subscribers: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			if resp.JSON200 != nil {
				for _, uid := range *resp.JSON200 {
					if _, err := fmt.Fprintln(out, uid); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
}

// NewTasksChatSubscribersUpdateCmd returns the "tasks chat-subscribers update" command.
func NewTasksChatSubscribersUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var userIDs string
	c := &cobra.Command{
		Use:   "update [task-id]",
		Short: "Update task chat subscribers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if userIDs == "" {
				return fmt.Errorf("user-ids is required (--user-ids id1,id2,...)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			parts := strings.Split(userIDs, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			body := client.TaskControllerUpdateChatSubscribersJSONRequestBody{Content: &parts}
			resp, err := api.TaskControllerUpdateChatSubscribersWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update chat subscribers: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update chat subscribers: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			_, err = fmt.Fprintf(out, "Chat subscribers updated for task %s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&userIDs, "user-ids", "", "comma-separated list of user IDs")
	_ = c.MarkFlagRequired("user-ids")
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
	c.AddCommand(NewTasksCreateCmd(resolvePath, outputJSON))
	c.AddCommand(NewTasksUpdateCmd(resolvePath, outputJSON))
	chatSubs := &cobra.Command{Use: "chat-subscribers", Short: "Task chat subscribers"}
	chatSubs.AddCommand(NewTasksChatSubscribersGetCmd(resolvePath, outputJSON))
	chatSubs.AddCommand(NewTasksChatSubscribersUpdateCmd(resolvePath, outputJSON))
	c.AddCommand(chatSubs)
	return c
}
