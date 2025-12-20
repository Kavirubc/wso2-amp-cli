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

//orgListCmd

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
}
