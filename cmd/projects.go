package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		name, _ := cmd.Flags().GetString("name")
		displayName, _ := cmd.Flags().GetString("display-name")
		description, _ := cmd.Flags().GetString("description")
		pipeline, _ := cmd.Flags().GetString("pipeline")
		output, _ := cmd.Flags().GetString("output")

		// Use default org from config if not provided
		if org == "" {
			org = config.GetDefaultOrg()
		}
		if org == "" {
			return fmt.Errorf("organization is required. Use --org flag or set default with: amp config set default_org <name>")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Interactive mode: prompt for missing required fields
		if displayName == "" {
			fmt.Println(ui.TitleStyle.Render("📁 Create New Project"))
			fmt.Println()
			fmt.Print("? Display name: ")
			input, _ := reader.ReadString('\n')
			displayName = strings.TrimSpace(input)
			if displayName == "" {
				return fmt.Errorf("display name is required")
			}
		}

		// Prompt for description if not provided
		if description == "" && !cmd.Flags().Changed("description") {
			fmt.Print("? Description (optional): ")
			input, _ := reader.ReadString('\n')
			description = strings.TrimSpace(input)
		}

		// Prompt for pipeline if not provided
		if pipeline == "" {
			pipelines, err := client.ListDeploymentPipelines(org)
			if err != nil {
				return fmt.Errorf("failed to fetch pipelines: %w", err)
			}
			if len(pipelines) == 0 {
				return fmt.Errorf("no deployment pipelines available in organization '%s'", org)
			}

			// Show pipeline options
			fmt.Println("? Select deployment pipeline:")
			for i, p := range pipelines {
				if p.DisplayName != "" {
					fmt.Printf("  %d. %s (%s)\n", i+1, p.Name, p.DisplayName)
				} else {
					fmt.Printf("  %d. %s\n", i+1, p.Name)
				}
			}
			fmt.Printf("Enter selection [1]: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			// Default to first option
			selection := 1
			if input != "" {
				sel, err := strconv.Atoi(input)
				if err != nil || sel < 1 || sel > len(pipelines) {
					return fmt.Errorf("invalid selection: %s", input)
				}
				selection = sel
			}
			pipeline = pipelines[selection-1].Name
		}

		// Generate name from display name if not provided
		if name == "" {
			name = generateProjectName(displayName)
			if name == "" {
				return fmt.Errorf("could not generate valid project name from '%s'. Please provide a name with --name flag", displayName)
			}
		}

		// Build request
		var desc *string
		if description != "" {
			desc = &description
		}
		req := api.CreateProjectRequest{
			Name:               name,
			DisplayName:        displayName,
			Description:        desc,
			DeploymentPipeline: pipeline,
		}

		fmt.Println()
		fmt.Println("Creating project...")

		// Create project
		project, err := client.CreateProject(org, req)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(project)
		}

		// Success output
		fmt.Println()
		fmt.Println(ui.RenderSuccess("Project created successfully!"))
		fmt.Println()
		printProjectRow("Name:", project.Name)
		printProjectRow("Display Name:", project.DisplayName)
		printProjectRow("Organization:", project.OrgName)
		printProjectRow("Pipeline:", project.DeploymentPipeline)
		fmt.Println()

		return nil
	},
}

// generateProjectName converts display name to valid project name
func generateProjectName(displayName string) string {
	// Lowercase and replace spaces with hyphens
	name := strings.ToLower(displayName)
	name = strings.ReplaceAll(name, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	name = reg.ReplaceAllString(name, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Trim hyphens from start and end
	name = strings.Trim(name, "-")

	// Truncate to 63 characters (Kubernetes limit)
	if len(name) > 63 {
		name = name[:63]
		name = strings.TrimRight(name, "-")
	}

	return name
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
	projectsCmd.AddCommand(projectsCreateCmd)

	// Add --force flag to delete command
	projectsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	// Add flags for create command
	projectsCreateCmd.Flags().String("name", "", "Project name (auto-generated if not provided)")
	projectsCreateCmd.Flags().String("display-name", "", "Display name for the project")
	projectsCreateCmd.Flags().String("description", "", "Project description")
	projectsCreateCmd.Flags().String("pipeline", "", "Deployment pipeline name")
}
