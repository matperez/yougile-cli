package cmd

import (
	"context"
	"fmt"

	"github.com/angolovin/yougile-cli/internal/output"
	"github.com/angolovin/yougile-cli/pkg/client"
	"github.com/spf13/cobra"
)

// NewStickersStringListCmd returns the "stickers string list" command.
func NewStickersStringListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var includeDeleted bool
	c := &cobra.Command{
		Use:   "list",
		Short: "List string stickers",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.StringStickerControllerSearchParams{}
			if includeDeleted {
				params.IncludeDeleted = boolPtr(true)
			}
			resp, err := api.StringStickerControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list string stickers: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list string stickers: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list string stickers: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Name"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, s := range resp.JSON200.Content {
				rows = append(rows, []string{s.Id, s.Name})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().BoolVar(&includeDeleted, "include-deleted", false, "include deleted stickers")
	return c
}

// NewStickersStringGetCmd returns the "stickers string get" command.
func NewStickersStringGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get string sticker by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			resp, err := api.StringStickerControllerGetWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get string sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get string sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get string sticker: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			return output.PrintJSON(out, resp.JSON200)
		},
	}
}

// NewStickersSprintListCmd returns the "stickers sprint list" command.
func NewStickersSprintListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var includeDeleted bool
	c := &cobra.Command{
		Use:   "list",
		Short: "List sprint stickers",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			params := &client.SprintStickerControllerSearchParams{}
			if includeDeleted {
				params.IncludeDeleted = boolPtr(true)
			}
			resp, err := api.SprintStickerControllerSearchWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("list sprint stickers: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("list sprint stickers: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("list sprint stickers: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			headers := []string{"ID", "Name"}
			rows := make([][]string, 0, len(resp.JSON200.Content))
			for _, s := range resp.JSON200.Content {
				rows = append(rows, []string{s.Id, s.Name})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
	c.Flags().BoolVar(&includeDeleted, "include-deleted", false, "include deleted stickers")
	return c
}

// NewStickersSprintGetCmd returns the "stickers sprint get" command.
func NewStickersSprintGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get sprint sticker by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			resp, err := api.SprintStickerControllerGetStickerWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("get sprint sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get sprint sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get sprint sticker: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			return output.PrintJSON(out, resp.JSON200)
		},
	}
}

// NewStickersCmd returns the "stickers" parent command.
func NewStickersCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "stickers",
		Short: "String and sprint stickers",
	}
	stringCmd := &cobra.Command{Use: "string", Short: "String stickers"}
	stringCmd.AddCommand(NewStickersStringListCmd(resolvePath, outputJSON))
	stringCmd.AddCommand(NewStickersStringGetCmd(resolvePath, outputJSON))
	c.AddCommand(stringCmd)
	sprintCmd := &cobra.Command{Use: "sprint", Short: "Sprint stickers"}
	sprintCmd.AddCommand(NewStickersSprintListCmd(resolvePath, outputJSON))
	sprintCmd.AddCommand(NewStickersSprintGetCmd(resolvePath, outputJSON))
	c.AddCommand(sprintCmd)
	return c
}
