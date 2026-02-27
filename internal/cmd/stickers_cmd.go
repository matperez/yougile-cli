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

// NewStickersStringCreateCmd returns the "stickers string create" command.
func NewStickersStringCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a string sticker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required (--name)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.StringStickerControllerCreateJSONRequestBody{Name: name}
			resp, err := api.StringStickerControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create string sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create string sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "String sticker created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "sticker name")
	_ = c.MarkFlagRequired("name")
	return c
}

// NewStickersStringUpdateCmd returns the "stickers string update" command.
func NewStickersStringUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a string sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.StringStickerControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("name") {
				body.Name = &name
			}
			resp, err := api.StringStickerControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update string sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update string sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "String sticker updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&name, "name", "", "sticker name")
	return c
}

// NewStickersStringStatesListCmd returns the "stickers string states list" command (states from sticker get).
func NewStickersStringStatesListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "list [sticker-id]",
		Short: "List states of a string sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID := args[0]
			resp, err := api.StringStickerControllerGetWithResponse(context.Background(), stickerID)
			if err != nil {
				return fmt.Errorf("get string sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
				return fmt.Errorf("get string sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			states := resp.JSON200.States
			if states == nil {
				states = &[]client.StringStickerStateDto{}
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, states)
			}
			headers := []string{"ID", "Name"}
			rows := make([][]string, 0, len(*states))
			for _, s := range *states {
				rows = append(rows, []string{s.Id, s.Name})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
}

// NewStickersStringStatesGetCmd returns the "stickers string states get" command.
func NewStickersStringStatesGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [sticker-id] [state-id]",
		Short: "Get a string sticker state by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID, stateID := args[0], args[1]
			resp, err := api.StringStickerStateControllerGetWithResponse(context.Background(), stickerID, stateID, nil)
			if err != nil {
				return fmt.Errorf("get string sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get string sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get string sticker state: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			return output.PrintJSON(out, resp.JSON200)
		},
	}
}

// NewStickersStringStatesCreateCmd returns the "stickers string states create" command.
func NewStickersStringStatesCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "create [sticker-id]",
		Short: "Create a state for a string sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required (--name)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID := args[0]
			body := client.StringStickerStateControllerCreateJSONRequestBody{Name: name}
			resp, err := api.StringStickerStateControllerCreateWithResponse(context.Background(), stickerID, body)
			if err != nil {
				return fmt.Errorf("create string sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create string sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "String sticker state created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "state name")
	_ = c.MarkFlagRequired("name")
	return c
}

// NewStickersStringStatesUpdateCmd returns the "stickers string states update" command.
func NewStickersStringStatesUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "update [sticker-id] [state-id]",
		Short: "Update a string sticker state",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID, stateID := args[0], args[1]
			body := client.StringStickerStateControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("name") {
				body.Name = &name
			}
			resp, err := api.StringStickerStateControllerUpdateWithResponse(context.Background(), stickerID, stateID, body)
			if err != nil {
				return fmt.Errorf("update string sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update string sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "String sticker state updated: id=%s\n", stateID)
			return err
		},
	}
	c.Flags().StringVar(&name, "name", "", "state name")
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

// NewStickersSprintCreateCmd returns the "stickers sprint create" command.
func NewStickersSprintCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a sprint sticker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required (--name)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			body := client.SprintStickerControllerCreateJSONRequestBody{Name: name}
			resp, err := api.SprintStickerControllerCreateWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("create sprint sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create sprint sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Sprint sticker created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "sticker name")
	_ = c.MarkFlagRequired("name")
	return c
}

// NewStickersSprintUpdateCmd returns the "stickers sprint update" command.
func NewStickersSprintUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "update [id]",
		Short: "Update a sprint sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			id := args[0]
			body := client.SprintStickerControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("name") {
				body.Name = &name
			}
			resp, err := api.SprintStickerControllerUpdateWithResponse(context.Background(), id, body)
			if err != nil {
				return fmt.Errorf("update sprint sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update sprint sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Sprint sticker updated: id=%s\n", id)
			return err
		},
	}
	c.Flags().StringVar(&name, "name", "", "sticker name")
	return c
}

// NewStickersSprintStatesListCmd returns the "stickers sprint states list" command.
func NewStickersSprintStatesListCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "list [sticker-id]",
		Short: "List states of a sprint sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID := args[0]
			resp, err := api.SprintStickerControllerGetStickerWithResponse(context.Background(), stickerID)
			if err != nil {
				return fmt.Errorf("get sprint sticker: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 || resp.JSON200 == nil {
				return fmt.Errorf("get sprint sticker: HTTP %s", resp.HTTPResponse.Status)
			}
			states := resp.JSON200.States
			if states == nil {
				states = &[]client.SprintStickerStateDto{}
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, states)
			}
			headers := []string{"ID", "Name"}
			rows := make([][]string, 0, len(*states))
			for _, s := range *states {
				rows = append(rows, []string{s.Id, s.Name})
			}
			return output.PrintTable(out, headers, rows)
		},
	}
}

// NewStickersSprintStatesGetCmd returns the "stickers sprint states get" command.
func NewStickersSprintStatesGetCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	return &cobra.Command{
		Use:   "get [sticker-id] [state-id]",
		Short: "Get a sprint sticker state by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID, stateID := args[0], args[1]
			resp, err := api.SprintStickerStateControllerGetWithResponse(context.Background(), stickerID, stateID, nil)
			if err != nil {
				return fmt.Errorf("get sprint sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("get sprint sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("get sprint sticker state: empty response")
			}
			out := cmd.OutOrStdout()
			if outputJSON() {
				return output.PrintJSON(out, resp.JSON200)
			}
			return output.PrintJSON(out, resp.JSON200)
		},
	}
}

// NewStickersSprintStatesCreateCmd returns the "stickers sprint states create" command.
func NewStickersSprintStatesCreateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "create [sticker-id]",
		Short: "Create a state for a sprint sticker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required (--name)")
			}
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID := args[0]
			body := client.SprintStickerStateControllerCreateJSONRequestBody{Name: name}
			resp, err := api.SprintStickerStateControllerCreateWithResponse(context.Background(), stickerID, body)
			if err != nil {
				return fmt.Errorf("create sprint sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 201 {
				return fmt.Errorf("create sprint sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON201 != nil {
				return output.PrintJSON(out, resp.JSON201)
			}
			if resp.JSON201 != nil {
				_, err = fmt.Fprintf(out, "Sprint sticker state created: id=%s\n", resp.JSON201.Id)
				return err
			}
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "state name")
	_ = c.MarkFlagRequired("name")
	return c
}

// NewStickersSprintStatesUpdateCmd returns the "stickers sprint states update" command.
func NewStickersSprintStatesUpdateCmd(resolvePath func() (string, error), outputJSON func() bool) *cobra.Command {
	var name string
	c := &cobra.Command{
		Use:   "update [sticker-id] [state-id]",
		Short: "Update a sprint sticker state",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, api, err := loadConfigAndClient(resolvePath)
			if err != nil {
				return err
			}
			stickerID, stateID := args[0], args[1]
			body := client.SprintStickerStateControllerUpdateJSONRequestBody{}
			if cmd.Flags().Changed("name") {
				body.Name = &name
			}
			resp, err := api.SprintStickerStateControllerUpdateWithResponse(context.Background(), stickerID, stateID, body)
			if err != nil {
				return fmt.Errorf("update sprint sticker state: %w", err)
			}
			if resp.HTTPResponse.StatusCode != 200 {
				return fmt.Errorf("update sprint sticker state: HTTP %s", resp.HTTPResponse.Status)
			}
			out := cmd.OutOrStdout()
			if outputJSON() && resp.JSON200 != nil {
				return output.PrintJSON(out, resp.JSON200)
			}
			_, err = fmt.Fprintf(out, "Sprint sticker state updated: id=%s\n", stateID)
			return err
		},
	}
	c.Flags().StringVar(&name, "name", "", "state name")
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
	stringCmd.AddCommand(NewStickersStringCreateCmd(resolvePath, outputJSON))
	stringCmd.AddCommand(NewStickersStringUpdateCmd(resolvePath, outputJSON))
	stringStates := &cobra.Command{Use: "states", Short: "String sticker states"}
	stringStates.AddCommand(NewStickersStringStatesListCmd(resolvePath, outputJSON))
	stringStates.AddCommand(NewStickersStringStatesGetCmd(resolvePath, outputJSON))
	stringStates.AddCommand(NewStickersStringStatesCreateCmd(resolvePath, outputJSON))
	stringStates.AddCommand(NewStickersStringStatesUpdateCmd(resolvePath, outputJSON))
	stringCmd.AddCommand(stringStates)
	c.AddCommand(stringCmd)
	sprintCmd := &cobra.Command{Use: "sprint", Short: "Sprint stickers"}
	sprintCmd.AddCommand(NewStickersSprintListCmd(resolvePath, outputJSON))
	sprintCmd.AddCommand(NewStickersSprintGetCmd(resolvePath, outputJSON))
	sprintCmd.AddCommand(NewStickersSprintCreateCmd(resolvePath, outputJSON))
	sprintCmd.AddCommand(NewStickersSprintUpdateCmd(resolvePath, outputJSON))
	sprintStates := &cobra.Command{Use: "states", Short: "Sprint sticker states"}
	sprintStates.AddCommand(NewStickersSprintStatesListCmd(resolvePath, outputJSON))
	sprintStates.AddCommand(NewStickersSprintStatesGetCmd(resolvePath, outputJSON))
	sprintStates.AddCommand(NewStickersSprintStatesCreateCmd(resolvePath, outputJSON))
	sprintStates.AddCommand(NewStickersSprintStatesUpdateCmd(resolvePath, outputJSON))
	sprintCmd.AddCommand(sprintStates)
	c.AddCommand(sprintCmd)
	return c
}
