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

var dataplanesCmd = &cobra.Command{
	Use:     "dataplanes",
	Aliases: []string{"dp"},
	Short:   "Manage data planes",
	Long:    `Commands for listing and viewing data planes in an organization.`,
}

var dataplanesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all data planes in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		org, _ := cmd.Flags().GetString("org")
		output, _ := cmd.Flags().GetString("output")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

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

		// Build pagination options
		opts := api.ListOptions{Limit: limit, Offset: offset}

		dataplanes, total, err := client.ListDataPlanes(org, opts)
		if err != nil {
			return fmt.Errorf("failed to list data planes: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(dataplanes)
		}

		if len(dataplanes) == 0 {
			fmt.Println(ui.RenderWarning("No data planes found."))
			return nil
		}

		headers := []string{"NAME", "DISPLAY NAME", "DESCRIPTION"}
		rows := make([][]string, len(dataplanes))
		for i, dp := range dataplanes {
			// Truncate description if too long
			desc := dp.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			rows[i] = []string{
				dp.Name,
				valueOrDefault(dp.DisplayName, "-"),
				valueOrDefault(desc, "-"),
			}
		}

		title := fmt.Sprintf("🖥️  Data Planes in %s", org)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))
		fmt.Println(ui.RenderPaginationInfo(offset, limit, total))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dataplanesCmd)
	dataplanesCmd.AddCommand(dataplanesListCmd)
}
