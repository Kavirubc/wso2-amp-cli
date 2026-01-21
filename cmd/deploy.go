package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an agent to an environment",
	Long:  `Deploy an agent to a target environment using a specific build image.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
		imageID, _ := cmd.Flags().GetString("image")
		envVars, _ := cmd.Flags().GetStringSlice("env")
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
		if imageID == "" {
			return fmt.Errorf("image ID is required. Use --image flag")
		}

		// Parse environment variables (format: KEY=VALUE)
		var envList []api.EnvironmentVariable
		for _, ev := range envVars {
			parts := strings.SplitN(ev, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid environment variable format: %s (expected KEY=VALUE)", ev)
			}
			envList = append(envList, api.EnvironmentVariable{
				Key:   parts[0],
				Value: parts[1],
			})
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Build deploy request
		req := api.DeployAgentRequest{
			ImageId: imageID,
			Env:     envList,
		}

		// Execute deployment
		err := client.DeployAgent(org, project, agent, req)
		if err != nil {
			return fmt.Errorf("failed to deploy agent: %w", err)
		}

		// JSON output
		if output == "json" {
			result := map[string]string{
				"status":  "success",
				"agent":   agent,
				"imageId": imageID,
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(result)
		}

		// Success message
		fmt.Println(ui.RenderSuccess("Deployment triggered successfully!"))
		fmt.Println()
		printDeployRow("Agent:", agent)
		printDeployRow("Image:", imageID)
		if len(envList) > 0 {
			printDeployRow("Environment:", fmt.Sprintf("%d variable(s) set", len(envList)))
		}
		fmt.Println()

		// Show next steps
		fmt.Println(ui.SubtitleStyle.Render("Next steps:"))
		fmt.Printf("  • View deployments: amp deployments list --agent %s\n", agent)
		fmt.Printf("  • View endpoints: amp deployments endpoints --agent %s\n", agent)
		fmt.Println()

		return nil
	},
}

// printDeployRow prints a styled key-value row for deployment details
func printDeployRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Required flags
	deployCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	deployCmd.Flags().StringP("image", "i", "", "Build image ID (required)")

	// Optional flags
	deployCmd.Flags().StringSliceP("env", "e", nil, "Environment variables (KEY=VALUE, can be specified multiple times)")
}
