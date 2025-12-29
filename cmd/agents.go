package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
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
			fmt.Println(ui.RenderWarning("No agents found."))
			return nil
		}

		// Build table data
		headers := []string{"NAME", "DISPLAY NAME", "STATUS"}
		rows := make([][]string, len(agents))
		for i, agent := range agents {
			rows[i] = []string{
				agent.Name,
				agent.DisplayName,
				ui.StatusCell(agent.Status),
			}
		}

		// Render styled table
		title := fmt.Sprintf("%s Agents in %s/%s", ui.IconAgent, org, project)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))

		return nil
	},
}

var agentsGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]

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

		// Fetch agent from API
		agent, err := client.GetAgent(org, project, agentName)
		if err != nil {
			return fmt.Errorf("failed to get agent: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(agent)
		}

		// Pretty print agent details
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Agent: %s", ui.IconAgent, agent.Name)))
		fmt.Println()
		printAgentRow("Name:", agent.Name)
		printAgentRow("Display Name:", agent.DisplayName)
		printAgentRow("Description:", valueOrDefault(agent.Description, "(none)"))
		printAgentRow("Project:", agent.ProjectName)
		printAgentRow("Status:", ui.StatusCell(agent.Status))
		printAgentRow("Created At:", agent.CreatedAt.Format("2006-01-02 15:04:05"))

		// Show provisioning details if available
		if agent.Provisioning != nil {
			fmt.Println()
			fmt.Println(ui.SubtitleStyle.Render("  Provisioning:"))
			printAgentRow("    Type:", agent.Provisioning.Type)
			if agent.Provisioning.Repository != nil {
				printAgentRow("    URL:", agent.Provisioning.Repository.URL)
				printAgentRow("    Branch:", agent.Provisioning.Repository.Branch)
				if agent.Provisioning.Repository.AppPath != "" {
					printAgentRow("    App Path:", agent.Provisioning.Repository.AppPath)
				}
			}
		}

		// Show agent type if available
		if agent.AgentType != nil {
			fmt.Println()
			fmt.Println(ui.SubtitleStyle.Render("  Type:"))
			printAgentRow("    Type:", agent.AgentType.Type)
			if agent.AgentType.SubType != "" {
				printAgentRow("    SubType:", agent.AgentType.SubType)
			}
		}

		// Show runtime configs if available
		if agent.RuntimeConfigs != nil {
			fmt.Println()
			fmt.Println(ui.SubtitleStyle.Render("  Runtime:"))
			if agent.RuntimeConfigs.Language != "" {
				printAgentRow("    Language:", agent.RuntimeConfigs.Language)
			}
			if agent.RuntimeConfigs.LanguageVersion != "" {
				printAgentRow("    Version:", agent.RuntimeConfigs.LanguageVersion)
			}
			if agent.RuntimeConfigs.RunCommand != "" {
				printAgentRow("    Run Command:", agent.RuntimeConfigs.RunCommand)
			}
		}

		fmt.Println()
		return nil
	},
}

var agentsDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		force, _ := cmd.Flags().GetBool("force")

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

		// Confirm deletion unless --force is used
		if !force {
			fmt.Printf("Are you sure you want to delete agent '%s' in project '%s'? [y/N]: ", agentName, project)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println(ui.RenderWarning("Deletion cancelled."))
				return nil
			}
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Delete agent
		err := client.DeleteAgent(org, project, agentName)
		if err != nil {
			return fmt.Errorf("failed to delete agent: %w", err)
		}

		fmt.Println(ui.RenderSuccess(fmt.Sprintf("Agent '%s' deleted successfully.", agentName)))
		return nil
	},
}

// printAgentRow prints a styled key-value row for agent details
func printAgentRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsGetCmd)
	agentsCmd.AddCommand(agentsDeleteCmd)

	// Add --force flag to delete command
	agentsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
