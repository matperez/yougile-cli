package main

import (
	"os"
	"path/filepath"

	"github.com/angolovin/yougile-cli/internal/cmd"
	"github.com/spf13/cobra"
)

const configDir = "yougile-cli"
const configFile = "config.yaml"

var (
	configPath string
	outputJSON bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (default: ~/.config/yougile-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output as JSON")

	rootCmd.AddCommand(cmd.NewConfigCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewAuthCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewCompanyCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewUsersCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewProjectsCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewBoardsCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewColumnsCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewTasksCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewDepartmentsCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewWebhooksCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewFilesCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewChatsCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewStickersCmd(ResolveConfigPath, OutputJSON))
	rootCmd.AddCommand(cmd.NewCrmCmd(ResolveConfigPath, OutputJSON))
}

var rootCmd = &cobra.Command{
	Use:   "yougile",
	Short: "YouGile CLI — project management and CRM",
	Long:  "CLI for YouGile: tasks, projects, boards, users, and more.",
}

// ResolveConfigPath returns the config file path: flag value if set, otherwise default under user config dir.
func ResolveConfigPath() (string, error) {
	if configPath != "" {
		return configPath, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configDir, configFile), nil
}

// OutputJSON returns whether --json was set.
func OutputJSON() bool {
	return outputJSON
}
