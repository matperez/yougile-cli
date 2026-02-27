package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewChatsListCmd returns the "chats list" command (group chats).
func NewChatsListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	var title string

	c := &cobra.Command{
		Use:   "list",
		Short: "List group chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.GroupChatControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			if title != "" {
				params.Title = strPtr(title)
			}
			resp, err := api.GroupChatControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list chats: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list chats: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list chats: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Title"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, ch := range resp.JSON200.Content {
				rows = append(rows, []string{ch.Id, ch.Title})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	c.Flags().StringVar(&title, "title", "", "filter by title")
	return c
}

// NewChatGetCmd returns the "chats get" command.
func NewChatGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get group chat by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			resp, err := api.GroupChatControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get chat: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get chat: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get chat: empty response")
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

// NewChatsMessagesListCmd returns the "chats messages list" command.
func NewChatsMessagesListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var limit, offset int
	c := &cobra.Command{
		Use:   "list [chat-id]",
		Short: "List messages in a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			chatID := args[0]
			params := &client.ChatMessageControllerSearchParams{}
			if limit > 0 {
				params.Limit = float32Ptr(float32(limit))
			}
			if offset > 0 {
				params.Offset = float32Ptr(float32(offset))
			}
			resp, err := api.ChatMessageControllerSearchWithResponse(context.Background(), chatID, params)
			if err != nil {
				return fmt.Errorf("list messages: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list messages: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list messages: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"Id", "FromUserId", "Text"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, m := range resp.JSON200.Content {
				rows = append(rows, []string{strconv.FormatFloat(float64(m.Id), 'f', 0, 32), m.FromUserId, m.Text})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().IntVar(&limit, "limit", 50, "max items to return")
	c.Flags().IntVar(&offset, "offset", 0, "offset for pagination")
	return c
}

// NewChatsMessagesSendCmd returns the "chats messages send" command.
func NewChatsMessagesSendCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var text string
	c := &cobra.Command{
		Use:   "send [chat-id]",
		Short: "Send a message to a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("message text is required (--text)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			chatID := args[0]
			body := client.ChatMessageControllerSendMessageJSONRequestBody{
				Label:   "",
				Text:    text,
				TextHtml: text,
			}
			resp, err := api.ChatMessageControllerSendMessageWithResponse(context.Background(), chatID, body)
			if err != nil {
				return fmt.Errorf("send message: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("send message: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Message id: %v\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&text, "text", "", "message text")
	_ = c.MarkFlagRequired("text")
	return c
}

// NewChatsCmd returns the "chats" parent command.
func NewChatsCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "chats",
		Short: "Group chats and messages",
	}
	c.AddCommand(NewChatsListCmd(resolvePath, outputJSON))
	c.AddCommand(NewChatGetCmd(resolvePath, outputJSON))
	msgs := &cobra.Command{Use: "messages", Short: "Chat messages"}
	msgs.AddCommand(NewChatsMessagesListCmd(resolvePath, outputJSON))
	msgs.AddCommand(NewChatsMessagesSendCmd(resolvePath, outputJSON))
	c.AddCommand(msgs)
	return c
}
