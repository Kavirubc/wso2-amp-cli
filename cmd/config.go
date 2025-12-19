package cmd

import (
	"fmt"

	"github.com/Kavirubc/amp-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Manage CLI configuration",
	Long: `Manage amp-cli configuration settings.

Configuration is stored in ~/.amp/config.yaml
Available keys:
  api_url         - Base URL of the API server
  api_key_header  - Header name for API key (default: X-API-Key)
  api_key         - Your API key value
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
		fmt.Printf("✅ Set %s = %s\n", key, value)
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
			fmt.Printf("%s: (not set)\n", key)
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
		return nil
	},
}
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📋 Current Configuration:")
		fmt.Println()
		fmt.Printf("  api_url:         %s\n", valueOrDefault(config.GetAPIURL(), "(not set)"))
		fmt.Printf("  api_key_header:  %s\n", valueOrDefault(config.GetAPIKeyHeader(), "(not set)"))
		fmt.Printf("  api_key:         %s\n", maskValue(config.GetAPIKeyValue()))
		fmt.Println()
		fmt.Printf("  default_org:     %s\n", valueOrDefault(config.GetDefaultOrg(), "(not set)"))
		fmt.Printf("  default_project: %s\n", valueOrDefault(config.GetDefaultProject(), "(not set)"))
		fmt.Println()
		fmt.Printf("Config file: %s\n", config.ConfigFile())
	},
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