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

var orgsCmd = &cobra.Command{
	Use:   "orgs",
	Short: "Manage Organizations",
	Long:  `Commands for listing and viewing organizations.`,
}

var orgsGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific organization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgName := args[0]
		output, _ := cmd.Flags().GetString("output")

		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		org, err := client.GetOrganization(orgName)
		if err != nil {
			return fmt.Errorf("failed to get organization: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(org)
		}

		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("🏢 Organization: %s", org.Name)))
		fmt.Println()
		printOrgRow("Name:", org.Name)
		printOrgRow("Display Name:", valueOrDefault(org.DisplayName, "(none)"))
		printOrgRow("Description:", valueOrDefault(org.Description, "(none)"))
		printOrgRow("Namespace:", valueOrDefault(org.Namespace, "(none)"))
		printOrgRow("Created At:", org.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		return nil
	},
}

func printOrgRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

var orgsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all organizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		//Get the output flag
		output, _ := cmd.Flags().GetString("output")

		//API Client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		//API Call

		orgs, err := client.ListOrganizations()

		if err != nil {
			return fmt.Errorf("failed to list organizations: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(orgs)
		}

		if len(orgs) == 0 {
			fmt.Println(ui.RenderWarning("No organizations found."))
			return nil
		}

		headers := []string{"NAME", "CREATED AT"}

		rows := make([][]string, len(orgs))
		for i, org := range orgs {
			rows[i] = []string{
				org.Name,
				org.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}
		fmt.Println(ui.RenderTableWithTitle("🏢 Organizations", headers, rows))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orgsCmd)
	orgsCmd.AddCommand(orgsListCmd)
	orgsCmd.AddCommand(orgsGetCmd)
}
