package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Kavirubc/amp-cli/internal/api"
	"github.com/Kavirubc/amp-cli/internal/config"
	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage agents",
	Long:  `Commands for listing, creating, and managing AI agents.`,
}

var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get org and project from flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		output, _ := cmd.Flags().GetString("output")

		// Use defaults from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}
		if project == "" {
			project = config.GetDefaultProject()
		}

		// Validate required fields
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}
		if project == "" {
			return fmt.Errorf("project is required. Use --project flag or set default with: amp config set default_project <name>")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch agents from API
		agents, err := client.ListAgents(org, project)
		if err != nil {
			return fmt.Errorf("failed to list agents: %w", err)
		}

		// Output based on format
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(agents)
		}

		// Table output
		if len(agents) == 0 {
			fmt.Println("No agents found.")
			return nil
		}

		fmt.Printf("📋 Agents in %s/%s:\n\n", org, project)
		fmt.Printf("%-20s %-25s %-10s %-10s\n", "NAME", "DISPLAY NAME", "STATUS", "LANGUAGE")
		fmt.Printf("%-20s %-25s %-10s %-10s\n", "----", "------------", "------", "--------")
		for _, agent := range agents {
			fmt.Printf("%-20s %-25s %-10s %-10s\n",
				agent.Name,
				agent.DisplayName,
				agent.Status,
				agent.Language,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	agentsCmd.AddCommand(agentsListCmd)
}
