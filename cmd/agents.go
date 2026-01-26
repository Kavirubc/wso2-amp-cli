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
	"github.com/Kavirubc/wso2-amp-cli/internal/util"
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
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

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

		// Build pagination options
		opts := api.ListOptions{Limit: limit, Offset: offset}

		// Fetch agents from API
		agents, total, err := client.ListAgents(org, project, opts)
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
		fmt.Println(ui.RenderPaginationInfo(offset, limit, total))

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
			fmt.Println()
			fmt.Println(ui.WarningStyle.Render("⚠️  You are about to delete an agent"))
			fmt.Printf("   Agent: %s\n", agentName)
			fmt.Printf("   Organization: %s\n", org)
			fmt.Printf("   Project: %s\n", project)
			fmt.Println()
			fmt.Println(ui.MutedStyle.Render("This action cannot be undone."))
			fmt.Print("Are you sure? [y/N]: ")
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

var agentsLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View runtime logs for a deployed agent",
	Long: `Fetch and display runtime logs from a deployed agent.

Examples:
  amp agents logs --agent myagent --env development
  amp agents logs --agent myagent --env dev --since 1h
  amp agents logs --agent myagent --env dev --level ERROR,WARN
  amp agents logs --agent myagent --env dev --search "connection failed"
  amp agents logs --agent myagent --env dev --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
		since, _ := cmd.Flags().GetString("since")
		level, _ := cmd.Flags().GetString("level")
		search, _ := cmd.Flags().GetString("search")
		limit, _ := cmd.Flags().GetInt("limit")
		sort, _ := cmd.Flags().GetString("sort")
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
		if envName == "" {
			return fmt.Errorf("environment name is required. Use --env flag")
		}

		// Validate limit range
		if limit < 1 || limit > 1000 {
			return fmt.Errorf("limit must be between 1 and 1000")
		}

		// Validate sort order
		sortOrder := strings.ToLower(strings.TrimSpace(sort))
		if sortOrder != "asc" && sortOrder != "desc" {
			return fmt.Errorf("invalid sort value %q: must be 'asc' or 'desc'", sort)
		}

		// Build log request
		req := api.RuntimeLogRequest{
			EnvironmentName: envName,
			Limit:           limit,
			SortOrder:       sortOrder,
		}

		// Parse --since flag into start time
		if since != "" {
			startTime, err := util.ParseSinceDuration(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}
			req.StartTime = startTime.Format(time.RFC3339)
			req.EndTime = time.Now().Format(time.RFC3339)
		}

		// Parse log levels
		if level != "" {
			levels := strings.Split(level, ",")
			for i, l := range levels {
				levels[i] = strings.TrimSpace(strings.ToUpper(l))
			}
			req.LogLevels = levels
		}

		// Set search phrase
		if search != "" {
			req.SearchPhrase = search
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch logs from API
		logs, err := client.GetAgentRuntimeLogs(org, project, agentName, req)
		if err != nil {
			return fmt.Errorf("failed to get runtime logs: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(logs)
		}

		// Display logs
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Runtime Logs: %s (%s)", ui.IconAgent, agentName, envName)))
		fmt.Println()

		if len(logs.Logs) == 0 {
			fmt.Println(ui.RenderWarning("No logs found for the specified criteria."))
			return nil
		}

		// Print each log entry with timestamp and level
		for _, entry := range logs.Logs {
			timestamp := ui.FormatLogTimestamp(entry.Timestamp)
			levelPrefix := ui.FormatLogLevel(entry.LogLevel)
			if levelPrefix != "" {
				fmt.Printf("[%s] %s %s\n", timestamp, levelPrefix, entry.Log)
			} else {
				fmt.Printf("[%s] %s\n", timestamp, entry.Log)
			}
		}

		fmt.Println()
		fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("Showing %d log entries", len(logs.Logs))))

		return nil
	},
}

var agentsMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "View resource metrics for a deployed agent",
	Long: `Fetch and display CPU and memory metrics for a deployed agent.

Examples:
  amp agents metrics --agent myagent --env development
  amp agents metrics --agent myagent --env dev --since 1h
  amp agents metrics --agent myagent --env dev --start "2025-01-20T13:00:00Z" --end "2025-01-20T14:00:00Z"
  amp agents metrics --agent myagent --env dev --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
		since, _ := cmd.Flags().GetString("since")
		startTime, _ := cmd.Flags().GetString("start")
		endTime, _ := cmd.Flags().GetString("end")
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
		if envName == "" {
			return fmt.Errorf("environment name is required. Use --env flag")
		}

		// Build metrics request
		req := api.MetricsFilterRequest{
			EnvironmentName: envName,
		}

		// Handle time range options
		if since != "" {
			// Parse --since flag
			start, err := util.ParseSinceDuration(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}
			req.StartTime = start.Format(time.RFC3339)
			req.EndTime = time.Now().Format(time.RFC3339)
		} else if startTime != "" || endTime != "" {
			// Use explicit start/end times
			if startTime == "" || endTime == "" {
				return fmt.Errorf("both --start and --end must be provided together")
			}
			// Validate RFC3339 format
			if _, err := time.Parse(time.RFC3339, startTime); err != nil {
				return fmt.Errorf("invalid --start time format. Use RFC3339 format (e.g., 2025-01-20T13:00:00Z)")
			}
			if _, err := time.Parse(time.RFC3339, endTime); err != nil {
				return fmt.Errorf("invalid --end time format. Use RFC3339 format (e.g., 2025-01-20T14:00:00Z)")
			}
			req.StartTime = startTime
			req.EndTime = endTime
		} else {
			// Default to last 1 hour
			req.StartTime = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
			req.EndTime = time.Now().Format(time.RFC3339)
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch metrics from API
		metrics, err := client.GetAgentMetrics(org, project, agentName, req)
		if err != nil {
			return fmt.Errorf("failed to get metrics: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(metrics)
		}

		// Check if there's any data
		if !ui.HasMetricsData(metrics) {
			fmt.Println(ui.RenderWarning("No metrics data found for the specified criteria."))
			return nil
		}

		// Display metrics
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Resource Metrics for %s (%s)", ui.IconMetrics, agentName, envName)))
		fmt.Println()

		// Show time range
		startDisplay := ui.FormatMetricTimestamp(req.StartTime)
		endDisplay := ui.FormatMetricTimestamp(req.EndTime)
		fmt.Printf("  %s  %s - %s\n", ui.KeyStyle.Render("Time Range:"), startDisplay, endDisplay)
		fmt.Println()

		// CPU Usage table
		if len(metrics.CpuUsage) > 0 || len(metrics.CpuRequests) > 0 || len(metrics.CpuLimits) > 0 {
			fmt.Println(ui.SectionStyle.Render("  CPU Usage:"))
			headers, rows := ui.BuildCPUMetricsTable(metrics.CpuUsage, metrics.CpuRequests, metrics.CpuLimits)
			if len(rows) > 0 {
				fmt.Println(ui.RenderTable(headers, rows))
			}
			fmt.Println()
		}

		// Memory Usage table
		if len(metrics.Memory) > 0 || len(metrics.MemoryRequests) > 0 || len(metrics.MemoryLimits) > 0 {
			fmt.Println(ui.SectionStyle.Render("  Memory Usage:"))
			headers, rows := ui.BuildMemoryMetricsTable(metrics.Memory, metrics.MemoryRequests, metrics.MemoryLimits)
			if len(rows) > 0 {
				fmt.Println(ui.RenderTable(headers, rows))
			}
			fmt.Println()
		}

		return nil
	},
}

var agentsConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "View environment variables configured for a deployed agent",
	Long: `Fetch and display environment variables configured for an agent in a specific environment.

Examples:
  amp agents config --agent myagent --env development
  amp agents config --agent myagent --env dev --output json
  amp agents config --agent myagent --env dev --show-secrets`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
		output, _ := cmd.Flags().GetString("output")
		showSecrets, _ := cmd.Flags().GetBool("show-secrets")

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
		if envName == "" {
			return fmt.Errorf("environment name is required. Use --env flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch configuration from API
		configResp, err := client.GetAgentConfigurations(org, project, agentName, envName)
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(configResp)
		}

		// Check if there are any configurations
		if len(configResp.Configurations) == 0 {
			fmt.Println(ui.RenderWarning("No environment variables configured for this agent."))
			return nil
		}

		// Display configurations
		title := fmt.Sprintf("%s Environment Variables for %s (%s)", ui.IconAgent, agentName, envName)
		fmt.Println(ui.TitleStyle.Render(title))
		fmt.Println()

		// Build table data
		headers := []string{"KEY", "VALUE"}
		rows := make([][]string, len(configResp.Configurations))
		for i, cfg := range configResp.Configurations {
			value := cfg.Value
			// Mask sensitive values unless --show-secrets is specified
			if !showSecrets && isSensitiveKey(cfg.Key) {
				value = maskSensitiveValue(cfg.Value)
			}
			rows[i] = []string{cfg.Key, value}
		}

		// Render table
		fmt.Println(ui.RenderTable(headers, rows))
		fmt.Println()

		// Show hint about --show-secrets if any values are masked
		if !showSecrets && hasAnySensitiveKey(configResp.Configurations) {
			fmt.Println(ui.MutedStyle.Render("  Tip: Use --show-secrets to reveal masked values"))
			fmt.Println()
		}

		return nil
	},
}

// isSensitiveKey checks if a key name suggests it contains sensitive data
func isSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	sensitivePatterns := []string{
		"secret", "password", "passwd", "pwd", "token", "api_key", "apikey",
		"auth", "credential", "private", "key", "cert", "certificate",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerKey, pattern) {
			return true
		}
	}
	return false
}

// maskSensitiveValue returns a masked version of the value, showing first 2 and last 2 chars
func maskSensitiveValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// hasAnySensitiveKey checks if any configuration has a sensitive key
func hasAnySensitiveKey(configs []api.EnvironmentVariable) bool {
	for _, cfg := range configs {
		if isSensitiveKey(cfg.Key) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(agentsCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsGetCmd)
	agentsCmd.AddCommand(agentsTokenCmd)
	agentsCmd.AddCommand(agentsDeleteCmd)
	agentsCmd.AddCommand(agentsCreateCmd)
	agentsCmd.AddCommand(agentsLogsCmd)
	agentsCmd.AddCommand(agentsMetricsCmd)
	agentsCmd.AddCommand(agentsConfigCmd)

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

	// Add flags for logs command
	agentsLogsCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	agentsLogsCmd.Flags().StringP("env", "e", "", "Environment name (required)")
	agentsLogsCmd.Flags().String("since", "", "Show logs since duration (e.g., 1h, 24h, 7d)")
	agentsLogsCmd.Flags().String("level", "", "Filter by log levels (comma-separated: ERROR,WARN,INFO,DEBUG)")
	agentsLogsCmd.Flags().String("search", "", "Search phrase to filter logs")
	agentsLogsCmd.Flags().Int("limit", 100, "Maximum number of log entries to return")
	agentsLogsCmd.Flags().String("sort", "desc", "Sort order (asc/desc)")

	// Add flags for metrics command
	agentsMetricsCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	agentsMetricsCmd.Flags().StringP("env", "e", "", "Environment name (required)")
	agentsMetricsCmd.Flags().String("since", "", "Show metrics since duration (e.g., 1h, 24h, 7d)")
	agentsMetricsCmd.Flags().String("start", "", "Start time (RFC3339 format)")
	agentsMetricsCmd.Flags().String("end", "", "End time (RFC3339 format)")

	// Add flags for config command
	agentsConfigCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	agentsConfigCmd.Flags().StringP("env", "e", "", "Environment name (required)")
	agentsConfigCmd.Flags().Bool("show-secrets", false, "Show unmasked values for sensitive variables")
}
