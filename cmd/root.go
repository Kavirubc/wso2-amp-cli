package cmd

import (
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "amp",
	Short: "A simple CLI tool for managing WSO2 AI Agent Management Platform",
	Long: `amp-cli lets you manage agents, builds, and deployments 
from your terminal.

Examples:
  amp agents list
  amp builds trigger --agent my-agent
  amp deploy --agent my-agent --build latest`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize config BEFORE any command runs
	cobra.OnInitialize(initConfig)

	// Persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringP("org", "o", "", "Organization name")
	rootCmd.PersistentFlags().StringP("project", "p", "", "Project name")
	rootCmd.PersistentFlags().StringP("output", "", "table", "Output format (table|json)")
}

// initConfig is called before any command executes
func initConfig() {
	_ = config.Init()
}
