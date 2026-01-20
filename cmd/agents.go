package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

var agentsTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate a JWT token for an agent",
	Long: `Generate a JWT token that allows an agent to authenticate with AMP.

This command requests a signed token for the specified agent in a given
organization and project. The token can then be used by the agent to
connect to AMP and perform actions permitted by its configuration.

The generated token is sensitive credential material:
  - Treat it like a password or API key
  - Store it securely (e.g., in a secrets manager)
  - Do not commit it to version control
  - Prefer shorter expiration times where possible

Once generated, the token will only be displayed once.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		expiresIn, _ := cmd.Flags().GetString("expires-in")
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
		if agentName == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Build token request
		req := &api.TokenRequest{}
		if expiresIn != "" {
			req.ExpiresIn = expiresIn
		}

		// Generate token
		tokenResp, err := client.GenerateAgentToken(org, project, agentName, req)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(tokenResp)
		}

		// Format timestamps
		issuedAt := time.Unix(tokenResp.IssuedAt, 0)
		expiresAt := time.Unix(tokenResp.ExpiresAt, 0)

		// Pretty print token details
		fmt.Println(ui.TitleStyle.Render("🔐 Agent Token Generated"))
		fmt.Println()
		printAgentRow("Agent:", agentName)
		printAgentRow("Project:", project)
		printAgentRow("Token Type:", tokenResp.TokenType)
		printAgentRow("Issued At:", issuedAt.Format("2006-01-02 15:04:05"))
		printAgentRow("Expires At:", expiresAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
		fmt.Println(ui.SubtitleStyle.Render("  Token:"))
		fmt.Printf("  %s\n", tokenResp.Token)
		fmt.Println()
		fmt.Println(ui.RenderWarning("Store this token securely. It will not be shown again."))
		fmt.Println(ui.RenderInfo("Tip: Use --output json and redirect to a file for secure storage:"))
		fmt.Printf("  amp agents token --agent %s --output json > token.json\n", agentName)

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

var agentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")
		displayName, _ := cmd.Flags().GetString("display-name")
		description, _ := cmd.Flags().GetString("description")
		provisioning, _ := cmd.Flags().GetString("provisioning")
		repoURL, _ := cmd.Flags().GetString("repo-url")
		branch, _ := cmd.Flags().GetString("branch")
		appPath, _ := cmd.Flags().GetString("app-path")
		subtype, _ := cmd.Flags().GetString("subtype")
		language, _ := cmd.Flags().GetString("language")
		languageVersion, _ := cmd.Flags().GetString("language-version")
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

		// Interactive mode: prompt for missing required fields
		if displayName == "" {
			fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Create New Agent", ui.IconAgent)))
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

		// Prompt for provisioning type if not provided
		if provisioning == "" {
			provOptions := []string{"internal", "external"}
			provDescs := []string{"Platform-hosted agent", "Self-hosted agent"}
			fmt.Println("? Provisioning type:")
			for i, opt := range provOptions {
				fmt.Printf("  %d. %s - %s\n", i+1, opt, provDescs[i])
			}
			fmt.Print("Enter selection [1]: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			selection := 1
			if input != "" {
				sel, err := strconv.Atoi(input)
				if err != nil || sel < 1 || sel > len(provOptions) {
					return fmt.Errorf("invalid selection: %s", input)
				}
				selection = sel
			}
			provisioning = provOptions[selection-1]
		}

		// For internal provisioning, gather repository and runtime details
		var repoConfig *api.RepositoryConfig
		var runtimeConfig *api.RuntimeConfig

		if provisioning == "internal" {
			// Repository URL
			if repoURL == "" {
				fmt.Print("? Repository URL (https://github.com/owner/repo): ")
				input, _ := reader.ReadString('\n')
				repoURL = strings.TrimSpace(input)
				if repoURL == "" {
					return fmt.Errorf("repository URL is required for internal agents")
				}
			}

			// Branch
			if branch == "" {
				fmt.Print("? Branch [main]: ")
				input, _ := reader.ReadString('\n')
				branch = strings.TrimSpace(input)
				if branch == "" {
					branch = "main"
				}
			}

			// App path
			if appPath == "" {
				fmt.Print("? App path [/]: ")
				input, _ := reader.ReadString('\n')
				appPath = strings.TrimSpace(input)
				if appPath == "" {
					appPath = "/"
				}
			}

			repoConfig = &api.RepositoryConfig{
				URL:     repoURL,
				Branch:  branch,
				AppPath: appPath,
			}

			// Agent subtype
			if subtype == "" {
				subtypeOptions := []string{"chat-api", "custom-api"}
				subtypeDescs := []string{"Conversational chat agent", "Custom API agent"}
				fmt.Println("? Agent subtype:")
				for i, opt := range subtypeOptions {
					fmt.Printf("  %d. %s - %s\n", i+1, opt, subtypeDescs[i])
				}
				fmt.Print("Enter selection [1]: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				selection := 1
				if input != "" {
					sel, err := strconv.Atoi(input)
					if err != nil || sel < 1 || sel > len(subtypeOptions) {
						return fmt.Errorf("invalid selection: %s", input)
					}
					selection = sel
				}
				subtype = subtypeOptions[selection-1]
			}

			// Language selection
			if language == "" {
				langOptions := []string{"python", "nodejs", "java", "go", "ballerina"}
				fmt.Println("? Language:")
				for i, opt := range langOptions {
					fmt.Printf("  %d. %s\n", i+1, opt)
				}
				fmt.Print("Enter selection [1]: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				selection := 1
				if input != "" {
					sel, err := strconv.Atoi(input)
					if err != nil || sel < 1 || sel > len(langOptions) {
						return fmt.Errorf("invalid selection: %s", input)
					}
					selection = sel
				}
				language = langOptions[selection-1]
			}

			// Language version (not required for ballerina)
			if languageVersion == "" && language != "ballerina" {
				defaultVersion := getDefaultVersion(language)
				fmt.Printf("? Language version [%s]: ", defaultVersion)
				input, _ := reader.ReadString('\n')
				languageVersion = strings.TrimSpace(input)
				if languageVersion == "" {
					languageVersion = defaultVersion
				}
			}

			runtimeConfig = &api.RuntimeConfig{
				Language:        language,
				LanguageVersion: languageVersion,
			}
		}

		// Build InputInterface for internal agents
		var inputInterface *api.InputInterface
		if provisioning == "internal" {
			inputInterface = &api.InputInterface{
				Type:     "HTTP",
				Port:     8080,
				BasePath: "/",
			}

			// For custom-api, prompt for additional details
			if subtype == "custom-api" {
				// Port
				fmt.Print("? HTTP Port [8080]: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					port, err := strconv.Atoi(input)
					if err != nil || port < 1 || port > 65535 {
						return fmt.Errorf("invalid port: %s", input)
					}
					inputInterface.Port = port
				}

				// Base path
				fmt.Print("? Base path [/]: ")
				input, _ = reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					inputInterface.BasePath = input
				}

				// Schema path (required for custom-api)
				fmt.Print("? OpenAPI schema path (e.g., /openapi.yaml): ")
				input, _ = reader.ReadString('\n')
				schemaPath := strings.TrimSpace(input)
				if schemaPath == "" {
					return fmt.Errorf("schema path is required for custom-api agents")
				}
				inputInterface.Schema = &api.SchemaConfig{Path: schemaPath}
			}
		}

		// Generate name from display name if not provided
		if name == "" {
			name = generateAgentName(displayName)
			if name == "" {
				return fmt.Errorf("could not generate valid agent name from '%s'. Please provide a name with --name flag", displayName)
			}
		}

		// Build request
		req := api.CreateAgentRequest{
			Name:        name,
			DisplayName: displayName,
			Description: description,
			Provisioning: api.Provisioning{
				Type:       provisioning,
				Repository: repoConfig,
			},
			AgentType: api.AgentTypeInfo{
				Type:    "api",
				SubType: subtype,
			},
			RuntimeConfigs: runtimeConfig,
			InputInterface: inputInterface,
		}

		fmt.Println()
		fmt.Println("Creating agent...")

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Create agent
		agent, err := client.CreateAgent(org, project, req)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(agent)
		}

		// Success output
		fmt.Println()
		fmt.Println(ui.RenderSuccess("Agent created successfully!"))
		fmt.Println()
		printAgentRow("Name:", agent.Name)
		printAgentRow("Display Name:", agent.DisplayName)
		printAgentRow("Provisioning:", provisioning)
		if provisioning == "internal" && runtimeConfig != nil {
			langInfo := runtimeConfig.Language
			if runtimeConfig.LanguageVersion != "" {
				langInfo += " " + runtimeConfig.LanguageVersion
			}
			printAgentRow("Language:", langInfo)
		}
		fmt.Println()

		// Show next steps
		fmt.Println(ui.SubtitleStyle.Render("Next steps:"))
		fmt.Printf("  • Trigger a build: amp builds trigger --agent %s\n", agent.Name)
		fmt.Printf("  • View agent: amp agents get %s\n", agent.Name)
		fmt.Println()

		return nil
	},
}

// generateAgentName converts display name to valid agent name (max 25 chars)
func generateAgentName(displayName string) string {
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

	// Ensure starts with a letter (required by API)
	if len(name) > 0 && (name[0] < 'a' || name[0] > 'z') {
		name = "a" + name
	}

	// Truncate to 25 characters (API limit)
	if len(name) > 25 {
		name = name[:25]
		name = strings.TrimRight(name, "-")
	}

	return name
}

// getDefaultVersion returns the default version for a given language
func getDefaultVersion(language string) string {
	defaults := map[string]string{
		"python": "3.11",
		"nodejs": "20.x.x",
		"java":   "21",
		"go":     "1.x",
	}
	if v, ok := defaults[language]; ok {
		return v
	}
	return ""
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsGetCmd)
	agentsCmd.AddCommand(agentsTokenCmd)
	agentsCmd.AddCommand(agentsDeleteCmd)
	agentsCmd.AddCommand(agentsCreateCmd)

	// Add flags for token command
	agentsTokenCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	agentsTokenCmd.Flags().String("expires-in", "", "Token expiry duration (e.g., 720h for 30 days)")

	// Add --force flag to delete command
	agentsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	// Add flags for create command
	agentsCreateCmd.Flags().String("name", "", "Agent name (auto-generated if not provided)")
	agentsCreateCmd.Flags().String("display-name", "", "Display name for the agent")
	agentsCreateCmd.Flags().String("description", "", "Agent description")
	agentsCreateCmd.Flags().String("provisioning", "", "Provisioning type (internal/external)")
	agentsCreateCmd.Flags().String("repo-url", "", "Repository URL (for internal agents)")
	agentsCreateCmd.Flags().String("branch", "", "Git branch (default: main)")
	agentsCreateCmd.Flags().String("app-path", "", "App path in repository (default: /)")
	agentsCreateCmd.Flags().String("subtype", "", "Agent subtype (chat-api/custom-api)")
	agentsCreateCmd.Flags().String("language", "", "Programming language")
	agentsCreateCmd.Flags().String("language-version", "", "Language version")
}
