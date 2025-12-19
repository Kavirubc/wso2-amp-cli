package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage agents",
	Long:  `Commands for listing, creating, and managing AI agents.`,
}

var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag values
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")

		// For now, just print what we would do
		fmt.Println("📋 Listing agents...")
		fmt.Printf("   Organization: %s\n", org)
		fmt.Printf("   Project: %s\n", project)
		fmt.Println("\n(API call will go here later)")
	},
}

func init() {
	// Add agents command to root
	rootCmd.AddCommand(agentsCmd)

	// Add list subcommand to agents
	agentsCmd.AddCommand(agentsListCmd)
}
