package cmd

import (
	"fmt"

	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long: `Manage amp-cli configuration settings.

Configuration is stored in ~/.amp/config.yaml
Available keys:
  api_url         - Base URL of the API server (default: http://localhost:8080/api/v1)
  api_key_header  - Auth header name (default: Authorization)
  api_key         - Your JWT token (e.g., Bearer eyJ...)
  default_org     - Default organization name
  default_project - Default project name`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2), // Requires exactly 2 arguments
	Example: `  amp config set api_url http://localhost:8080
  amp config set api_key your-secret-key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		if err := config.Set(key, value); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println(ui.RenderSuccess(fmt.Sprintf("Set %s = %s", key, value)))
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := config.Get(key)
		if value == "" {
			fmt.Printf("%s  %s\n", ui.KeyStyle.Render(key+":"), ui.MutedStyle.Render("(not set)"))
		} else {
			fmt.Printf("%s  %s\n", ui.KeyStyle.Render(key+":"), ui.ValueStyle.Render(value))
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Current Configuration", ui.IconConfig)))
		fmt.Println()

		// API Settings
		fmt.Println(ui.SectionStyle.Render("API Settings"))
		printConfigRow("api_url", valueOrDefault(config.GetAPIURL(), "(not set)"), false)
		printConfigRow("api_key_header", valueOrDefault(config.GetAPIKeyHeader(), "(not set)"), false)
		printConfigRow("api_key", maskValue(config.GetAPIKeyValue()), true)
		fmt.Println()

		// Default Values
		fmt.Println(ui.SectionStyle.Render("Defaults"))
		printConfigRow("default_org", valueOrDefault(config.GetDefaultOrg(), "(not set)"), false)
		printConfigRow("default_project", valueOrDefault(config.GetDefaultProject(), "(not set)"), false)
		fmt.Println()

		// Config file location
		fmt.Printf("%s  %s\n", ui.MutedStyle.Render("Config file:"), ui.ValueStyle.Render(config.ConfigFile()))
	},
}

// printConfigRow prints a styled config key-value pair
func printConfigRow(key, value string, masked bool) {
	keyStr := ui.KeyStyle.Render(key + ":")
	var valStr string
	if value == "(not set)" {
		valStr = ui.MutedStyle.Render(value)
	} else if masked {
		valStr = ui.MaskedStyle.Render(value)
	} else {
		valStr = ui.ValueStyle.Render(value)
	}
	fmt.Printf("  %s  %s\n", keyStr, valStr)
}

// Helper functions
func valueOrDefault(value, defaultVal string) string {
	if value == "" {
		return defaultVal
	}
	return value
}

func maskValue(value string) string {
	if value == "" {
		return "(not set)"
	}
	if len(value) > 4 {
		return value[:4] + "****"
	}
	return "****"
}

func init() {
	// Add config command to root
	rootCmd.AddCommand(configCmd)

	// Add subcommands to config
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}
