package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

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

var orgsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new organization",
	Long: `Create a new organization in the platform.

You can provide the organization name as an argument or use the --name flag.
If neither is provided, you will be prompted to enter it interactively.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		output, _ := cmd.Flags().GetString("output")

		// Get name from args, flag, or prompt
		name, _ := cmd.Flags().GetString("name")
		if name == "" && len(args) > 0 {
			name = args[0]
		}

		// Interactive prompt if name not provided
		if name == "" {
			fmt.Println(ui.TitleStyle.Render("🏢 Create New Organization"))
			fmt.Println()
			fmt.Print("? Organization name: ")
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}

		// Validate name
		if name == "" {
			return fmt.Errorf("organization name is required")
		}
		name = generateOrgName(name)
		if name == "" {
			return fmt.Errorf("invalid organization name: must contain alphanumeric characters")
		}

		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		fmt.Println()
		fmt.Println("Creating organization...")

		req := api.CreateOrganizationRequest{Name: name}
		org, err := client.CreateOrganization(req)
		if err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}

		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(org)
		}

		fmt.Println()
		fmt.Println(ui.RenderSuccess("Organization created successfully!"))
		fmt.Println()
		printOrgRow("Name:", org.Name)
		printOrgRow("Display Name:", valueOrDefault(org.DisplayName, "(none)"))
		printOrgRow("Namespace:", valueOrDefault(org.Namespace, "(none)"))
		printOrgRow("Created At:", org.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
		fmt.Println(ui.MutedStyle.Render("Next steps:"))
		fmt.Printf("  • Set as default: amp config set default_org %s\n", org.Name)
		fmt.Printf("  • Create a project: amp projects create --org %s\n", org.Name)
		fmt.Println()

		return nil
	},
}

// generateOrgName sanitizes input to create a valid organization name
func generateOrgName(input string) string {
	name := strings.ToLower(strings.TrimSpace(input))
	name = strings.ReplaceAll(name, " ", "-")

	// Keep only alphanumeric and hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	name = reg.ReplaceAllString(name, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Trim hyphens from edges
	name = strings.Trim(name, "-")

	// Enforce max length (63 chars for Kubernetes compatibility)
	if len(name) > 63 {
		name = strings.TrimRight(name[:63], "-")
	}

	return name
}

func init() {
	rootCmd.AddCommand(orgsCmd)
	orgsCmd.AddCommand(orgsListCmd)
	orgsCmd.AddCommand(orgsGetCmd)
	orgsCmd.AddCommand(orgsCreateCmd)

	// Flags for create command
	orgsCreateCmd.Flags().String("name", "", "Organization name")
}
