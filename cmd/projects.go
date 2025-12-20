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

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `Commands for listing, viewing, and managing projects within an organization.`,
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		output, _ := cmd.Flags().GetString("output")

		// Use default from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}

		// Validate required fields
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch projects from API
		projects, err := client.ListProjects(org)
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		// Output based on format
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(projects)
		}

		// Table output
		if len(projects) == 0 {
			fmt.Println(ui.RenderWarning("No projects found."))
			return nil
		}

		// Build table data
		headers := []string{"NAME", "DISPLAY NAME", "CREATED AT"}
		rows := make([][]string, len(projects))
		for i, project := range projects {
			rows[i] = []string{
				project.Name,
				project.DisplayName,
				project.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		// Render styled table
		title := fmt.Sprintf("📁 Projects in %s", org)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))

		return nil
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Get org flag
		org, _ := cmd.Flags().GetString("org")
		output, _ := cmd.Flags().GetString("output")

		// Use default from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}

		// Validate required fields
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch project from API
		project, err := client.GetProject(org, projectName)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		// Output based on format
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(project)
		}

		// Pretty print project details
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("📁 Project: %s", project.Name)))
		fmt.Println()
		printProjectRow("Name:", project.Name)
		printProjectRow("Display Name:", project.DisplayName)
		printProjectRow("Description:", valueOrDefault(project.Description, "(none)"))
		printProjectRow("Organization:", project.OrgName)
		if project.DeploymentPipeline != "" {
			printProjectRow("Deployment Pipeline:", project.DeploymentPipeline)
		}
		printProjectRow("Created At:", project.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		return nil
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		force, _ := cmd.Flags().GetBool("force")

		// Use default from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}

		// Validate required fields
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}

		// Confirm deletion unless --force is used
		if !force {
			fmt.Printf("Are you sure you want to delete project '%s' in organization '%s'? [y/N]: ", projectName, org)
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

		// Delete project
		err := client.DeleteProject(org, projectName)
		if err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}

		fmt.Println(ui.RenderSuccess(fmt.Sprintf("Project '%s' deleted successfully.", projectName)))
		return nil
	},
}

// printProjectRow prints a styled key-value row for project details
func printProjectRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	// Add projects command to root
	rootCmd.AddCommand(projectsCmd)

	// Add subcommands to projects
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)

	// Add --force flag to delete command
	projectsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
