package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Kavirubc/wso2-amp-cli/internal/cli"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "amp",
	Short: "CLI for WSO2 AI Agent Management Platform",
	Long: `amp-cli lets you manage organizations, projects, and agents
from your terminal.

Examples:
  amp orgs list
  amp projects list --org default
  amp agents list --org default --project myproject
  amp config set default_org myorg`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show banner
		org := config.GetDefaultOrg()
		project := config.GetDefaultProject()
		fmt.Println(ui.RenderBanner(org, project))

		// Start interactive mode
		executor := cli.NewExecutor()
		model := ui.NewInteractiveModel(executor.Execute)

		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
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
