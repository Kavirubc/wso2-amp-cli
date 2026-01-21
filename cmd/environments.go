package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var environmentsCmd = &cobra.Command{
	Use:     "environments",
	Aliases: []string{"envs", "env"},
	Short:   "Manage environments",
	Long:    `Commands for listing and viewing environments in an organization.`,
}

var environmentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		org, _ := cmd.Flags().GetString("org")
		output, _ := cmd.Flags().GetString("output")

		// Use default org from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}

		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		environments, err := client.ListEnvironments(org)
		if err != nil {
			return fmt.Errorf("failed to list environments: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(environments)
		}

		if len(environments) == 0 {
			fmt.Println(ui.RenderWarning("No environments found."))
			return nil
		}

		headers := []string{"NAME", "DISPLAY NAME", "PRODUCTION", "CREATED AT"}
		rows := make([][]string, len(environments))
		for i, env := range environments {
			prodStatus := "No"
			if env.IsProduction {
				prodStatus = "Yes ★"
			}
			rows[i] = []string{
				env.Name,
				valueOrDefault(env.DisplayName, "-"),
				prodStatus,
				env.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		title := fmt.Sprintf("🌍 Environments in %s", org)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(environmentsCmd)
	environmentsCmd.AddCommand(environmentsListCmd)
}
