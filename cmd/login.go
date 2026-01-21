package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and configure the CLI",
	Long:  `Interactive setup wizard for configuring API connection and authentication.`,
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Get flags for non-interactive mode
	apiURL, _ := cmd.Flags().GetString("api-url")
	token, _ := cmd.Flags().GetString("token")

	fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s WSO2 AMP CLI Setup", ui.IconConfig)))
	fmt.Println()
	fmt.Println("Welcome! Let's configure your CLI.")
	fmt.Println()

	// Step 1: API URL
	if apiURL == "" {
		currentURL := config.GetAPIURL()
		if currentURL != "" && currentURL != "http://localhost:8080/api/v1" {
			fmt.Printf("? API Server URL [%s]: ", currentURL)
		} else {
			fmt.Print("? API Server URL: ")
		}
		input, _ := reader.ReadString('\n')
		apiURL = strings.TrimSpace(input)
		if apiURL == "" {
			apiURL = currentURL
		}
	}

	if apiURL == "" {
		return fmt.Errorf("API URL is required")
	}

	// Test connection
	fmt.Print("  Testing connection... ")
	if err := api.TestConnection(apiURL); err != nil {
		fmt.Println(ui.RenderError("Failed"))
		return fmt.Errorf("cannot connect to %s: %w", apiURL, err)
	}
	fmt.Println(ui.RenderSuccess("Connected"))

	// Step 2: Authentication token
	if token == "" {
		fmt.Println()
		fmt.Print("? API Token (paste your token): ")
		input, _ := reader.ReadString('\n')
		token = strings.TrimSpace(input)
	}

	if token == "" {
		return fmt.Errorf("API token is required")
	}

	// Determine auth header format
	authHeader := "Authorization"
	authValue := token
	if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
		authValue = "Bearer " + token
	}

	// Validate authentication BEFORE saving credentials
	fmt.Print("  Validating credentials... ")
	client := api.NewClient(apiURL, authHeader, authValue)
	orgs, err := client.ValidateAuth()
	if err != nil {
		fmt.Println(ui.RenderError("Failed"))
		return fmt.Errorf("authentication failed: %w", err)
	}
	fmt.Println(ui.RenderSuccess("Authenticated"))

	// Save credentials only after successful validation
	if err := config.Set(config.KeyAPIURL, apiURL); err != nil {
		return fmt.Errorf("failed to save API URL: %w", err)
	}
	if err := config.Set(config.KeyAPIKeyHeader, authHeader); err != nil {
		return fmt.Errorf("failed to save auth header: %w", err)
	}
	if err := config.Set(config.KeyAPIKeyValue, authValue); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Step 3: Select default organization
	fmt.Println()

	var defaultOrg string
	if len(orgs) == 0 {
		fmt.Println(ui.RenderWarning("No organizations found"))
	} else if len(orgs) == 1 {
		defaultOrg = orgs[0].Name
		fmt.Printf("  Default organization: %s\n", defaultOrg)
	} else {
		defaultOrg, err = selectOrganization(reader, orgs)
		if err != nil {
			return err
		}
	}

	if defaultOrg != "" {
		if err := config.Set(config.KeyDefaultOrg, defaultOrg); err != nil {
			return fmt.Errorf("failed to save default org: %w", err)
		}
	}

	// Step 4: Select default project
	if defaultOrg != "" {
		projects, err := client.ListProjects(defaultOrg)
		if err != nil {
			fmt.Println(ui.RenderWarning("Could not fetch projects"))
		} else if len(projects) > 0 {
			var defaultProj string
			if len(projects) == 1 {
				defaultProj = projects[0].Name
				fmt.Printf("  Default project: %s\n", defaultProj)
			} else {
				defaultProj, err = selectProject(reader, projects)
				if err != nil {
					return err
				}
			}

			if defaultProj != "" {
				if err := config.Set(config.KeyDefaultProj, defaultProj); err != nil {
					return fmt.Errorf("failed to save default project: %w", err)
				}
			}
		}
	}

	// Success message
	fmt.Println()
	fmt.Println(ui.RenderSuccess(fmt.Sprintf("Configuration saved to %s", config.ConfigFile())))
	fmt.Println()
	fmt.Println(ui.SubtitleStyle.Render("You're all set! Try these commands:"))
	fmt.Println("  • List agents: amp agents list")
	fmt.Println("  • View config: amp config show")
	fmt.Println("  • Get help: amp --help")
	fmt.Println()

	return nil
}

func selectOrganization(reader *bufio.Reader, orgs []api.OrganizationResponse) (string, error) {
	fmt.Println("? Default organization:")
	for i, org := range orgs {
		displayName := org.Name
		if org.DisplayName != "" {
			displayName = fmt.Sprintf("%s (%s)", org.DisplayName, org.Name)
		}
		fmt.Printf("  %d. %s\n", i+1, displayName)
	}
	fmt.Print("Enter selection [1]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	selection := 1
	if input != "" {
		sel, err := strconv.Atoi(input)
		if err != nil || sel < 1 || sel > len(orgs) {
			return "", fmt.Errorf("invalid selection: %s", input)
		}
		selection = sel
	}

	return orgs[selection-1].Name, nil
}

func selectProject(reader *bufio.Reader, projects []api.ProjectResponse) (string, error) {
	fmt.Println("? Default project:")
	for i, proj := range projects {
		displayName := proj.Name
		if proj.DisplayName != "" {
			displayName = fmt.Sprintf("%s (%s)", proj.DisplayName, proj.Name)
		}
		fmt.Printf("  %d. %s\n", i+1, displayName)
	}
	fmt.Print("Enter selection [1]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	selection := 1
	if input != "" {
		sel, err := strconv.Atoi(input)
		if err != nil || sel < 1 || sel > len(projects) {
			return "", fmt.Errorf("invalid selection: %s", input)
		}
		selection = sel
	}

	return projects[selection-1].Name, nil
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	Long:  `Remove stored authentication credentials from the CLI configuration.`,
	RunE:  runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Check if already logged out
	if !config.IsConfigured() {
		fmt.Println(ui.RenderInfo("No credentials stored. Already logged out."))
		return nil
	}

	// Get force flag
	force, _ := cmd.Flags().GetBool("force")

	// Confirm unless force flag is set
	if !force {
		fmt.Print("? Are you sure you want to log out? [y/N]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "y" && input != "yes" {
			fmt.Println("Logout cancelled.")
			return nil
		}
	}

	// Clear credentials
	if err := config.ClearCredentials(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Println(ui.RenderSuccess("Credentials removed"))
	fmt.Println()
	fmt.Println("To log in again, run: amp login")

	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	// Login flags for non-interactive mode
	loginCmd.Flags().String("api-url", "", "API server URL")
	loginCmd.Flags().String("token", "", "API token for authentication")

	// Logout flags
	logoutCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
