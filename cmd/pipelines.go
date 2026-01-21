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

var pipelinesCmd = &cobra.Command{
	Use:     "pipelines",
	Aliases: []string{"pipeline"},
	Short:   "Manage deployment pipelines",
	Long:    `Commands for listing and viewing deployment pipelines in an organization.`,
}

var pipelinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployment pipelines in an organization",
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

		pipelines, err := client.ListDeploymentPipelines(org)
		if err != nil {
			return fmt.Errorf("failed to list deployment pipelines: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(pipelines)
		}

		if len(pipelines) == 0 {
			fmt.Println(ui.RenderWarning("No deployment pipelines found."))
			return nil
		}

		headers := []string{"NAME", "DISPLAY NAME", "DESCRIPTION"}
		rows := make([][]string, len(pipelines))
		for i, p := range pipelines {
			// Truncate description if too long
			desc := p.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			rows[i] = []string{
				p.Name,
				valueOrDefault(p.DisplayName, "-"),
				valueOrDefault(desc, "-"),
			}
		}

		title := fmt.Sprintf("🔀 Deployment Pipelines in %s", org)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))
		return nil
	},
}

var pipelinesGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific deployment pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pipelineName := args[0]
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

		pipeline, err := client.GetDeploymentPipeline(org, pipelineName)
		if err != nil {
			return fmt.Errorf("failed to get deployment pipeline: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(pipeline)
		}

		// Display pipeline details
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("🔀 Deployment Pipeline: %s", pipeline.Name)))
		fmt.Println()
		printPipelineRow("Name:", pipeline.Name)
		printPipelineRow("Display Name:", valueOrDefault(pipeline.DisplayName, "(none)"))
		printPipelineRow("Description:", valueOrDefault(pipeline.Description, "(none)"))
		printPipelineRow("Organization:", pipeline.OrgName)
		printPipelineRow("Created At:", pipeline.CreatedAt.Format("2006-01-02 15:04:05"))

		// Display promotion paths if available
		if len(pipeline.PromotionPaths) > 0 {
			fmt.Println()
			fmt.Println(ui.SectionStyle.Render("Promotion Paths"))

			headers := []string{"FROM", "TO"}
			rows := make([][]string, 0)
			for _, path := range pipeline.PromotionPaths {
				targets := make([]string, len(path.TargetEnvironmentRefs))
				for i, t := range path.TargetEnvironmentRefs {
					targets[i] = t.Name
				}
				rows = append(rows, []string{
					path.SourceEnvironmentRef,
					strings.Join(targets, ", "),
				})
			}
			fmt.Println(ui.RenderTable(headers, rows))
		}

		fmt.Println()
		return nil
	},
}

func printPipelineRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.AddCommand(pipelinesListCmd)
	pipelinesCmd.AddCommand(pipelinesGetCmd)
}
