package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Manage agent deployments",
	Long:  `Commands for listing deployments and viewing endpoints for deployed agents.`,
}

var deploymentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments for an agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch deployments from API (map of env name to details)
		deployments, err := client.GetDeploymentsMap(org, project, agent)
		if err != nil {
			return fmt.Errorf("failed to list deployments: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(deployments)
		}

		// Table output
		if len(deployments) == 0 {
			fmt.Println(ui.RenderWarning("No deployments found."))
			return nil
		}

		// Sort environment names for consistent display
		envNames := make([]string, 0, len(deployments))
		for env := range deployments {
			envNames = append(envNames, env)
		}
		sort.Strings(envNames)

		// Build table data
		headers := []string{"ENVIRONMENT", "STATUS", "IMAGE", "LAST DEPLOYED", "ENDPOINTS"}
		rows := make([][]string, len(envNames))
		for i, env := range envNames {
			d := deployments[env]
			lastDeployed := "-"
			if d.LastDeployed != nil {
				lastDeployed = d.LastDeployed.Format("2006-01-02 15:04:05")
			}
			displayName := env
			if d.EnvironmentDisplayName != "" {
				displayName = d.EnvironmentDisplayName
			}
			rows[i] = []string{
				displayName,
				ui.StatusCell(d.Status),
				truncateImageID(d.ImageID),
				lastDeployed,
				fmt.Sprintf("%d", len(d.Endpoints)),
			}
		}

		// Render styled table
		title := fmt.Sprintf("%s Deployments for %s/%s/%s", ui.IconDeploy, org, project, agent)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))

		return nil
	},
}

var deploymentsEndpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "List endpoints for a deployed agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
		env, _ := cmd.Flags().GetString("env")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch endpoints from API
		endpoints, err := client.GetAgentEndpoints(org, project, agent, env)
		if err != nil {
			return fmt.Errorf("failed to get endpoints: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(endpoints)
		}

		// Table output
		if len(endpoints) == 0 {
			fmt.Println(ui.RenderWarning("No endpoints found."))
			return nil
		}

		// Build table data
		headers := []string{"NAME", "URL", "VISIBILITY"}
		rows := make([][]string, len(endpoints))
		for i, ep := range endpoints {
			rows[i] = []string{
				ep.EndpointName,
				ep.URL,
				ep.Visibility,
			}
		}

		// Render styled table
		title := fmt.Sprintf("%s Endpoints for %s/%s/%s", ui.IconEndpoint, org, project, agent)
		if env != "" {
			title = fmt.Sprintf("%s Endpoints for %s/%s/%s (%s)", ui.IconEndpoint, org, project, agent, env)
		}
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))

		return nil
	},
}

// truncateImageID shortens image ID for display (first 12 chars)
func truncateImageID(imageID string) string {
	if len(imageID) > 12 {
		return imageID[:12] + "..."
	}
	return imageID
}

func init() {
	rootCmd.AddCommand(deploymentsCmd)

	// Register subcommands
	deploymentsCmd.AddCommand(deploymentsListCmd)
	deploymentsCmd.AddCommand(deploymentsEndpointsCmd)

	// Add --agent flag to all subcommands (required)
	deploymentsListCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	deploymentsEndpointsCmd.Flags().StringP("agent", "a", "", "Agent name (required)")

	// Add --env flag to endpoints command (optional filter)
	deploymentsEndpointsCmd.Flags().StringP("env", "e", "", "Environment name (optional filter)")
}
