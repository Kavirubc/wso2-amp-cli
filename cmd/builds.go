package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
	"github.com/Kavirubc/wso2-amp-cli/internal/config"
	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

var buildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "Manage agent builds",
	Long:  `Commands for listing, viewing, triggering, and monitoring agent builds.`,
}

var buildsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all builds for an agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch builds from API
		builds, err := client.ListBuilds(org, project, agent)
		if err != nil {
			return fmt.Errorf("failed to list builds: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(builds)
		}

		// Table output
		if len(builds) == 0 {
			fmt.Println(ui.RenderWarning("No builds found."))
			return nil
		}

		// Build table data
		headers := []string{"NAME", "COMMIT", "STATUS", "STARTED", "DURATION"}
		rows := make([][]string, len(builds))
		for i, build := range builds {
			rows[i] = []string{
				build.Name,
				truncateCommit(build.CommitID),
				ui.StatusCell(build.Status),
				build.StartedAt.Format("2006-01-02 15:04:05"),
				formatDuration(build.StartedAt, build.EndedAt),
			}
		}

		// Render styled table
		title := fmt.Sprintf("%s Builds for %s/%s/%s", ui.IconBuild, org, project, agent)
		fmt.Println(ui.RenderTableWithTitle(title, headers, rows))

		return nil
	},
}

var buildsGetCmd = &cobra.Command{
	Use:   "get <build-name>",
	Short: "Get build details with steps",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		buildName := args[0]

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch build details from API
		build, err := client.GetBuild(org, project, agent, buildName)
		if err != nil {
			return fmt.Errorf("failed to get build: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(build)
		}

		// Pretty print build details
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Build: %s", ui.IconBuild, build.Name)))
		fmt.Println()
		printBuildRow("Name:", build.Name)
		printBuildRow("Agent:", build.AgentName)
		printBuildRow("Commit:", build.CommitID)
		printBuildRow("Branch:", valueOrDefault(build.Branch, "(unknown)"))
		printBuildRow("Status:", ui.StatusCell(build.Status))
		printBuildRow("Started:", build.StartedAt.Format("2006-01-02 15:04:05"))
		if build.EndedAt != nil {
			printBuildRow("Ended:", build.EndedAt.Format("2006-01-02 15:04:05"))
			printBuildRow("Duration:", formatDuration(build.StartedAt, build.EndedAt))
		} else {
			printBuildRow("Duration:", formatDuration(build.StartedAt, nil)+" (in progress)")
		}

		// Display progress percentage if available
		if build.Percent > 0 {
			printBuildRow("Progress:", fmt.Sprintf("%.1f%%", build.Percent))
		}

		// Display build steps if available
		if len(build.Steps) > 0 {
			fmt.Println()
			fmt.Println(ui.SubtitleStyle.Render("  Build Steps:"))
			for _, step := range build.Steps {
				stepStatus := ui.StatusCell(step.Status)
				fmt.Printf("    %s %s - %s\n", stepStatus, step.Type, step.Message)
			}
		}

		fmt.Println()
		return nil
	},
}

var buildsTriggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger a new build for an agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
		commit, _ := cmd.Flags().GetString("commit")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Trigger build
		build, err := client.TriggerBuild(org, project, agent, commit)
		if err != nil {
			return fmt.Errorf("failed to trigger build: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(build)
		}

		// Success message
		fmt.Println(ui.RenderSuccess("Build triggered successfully!"))
		fmt.Println()
		printBuildRow("Name:", build.Name)
		printBuildRow("Commit:", build.CommitID)
		printBuildRow("Status:", ui.StatusCell(build.Status))
		fmt.Println()

		// Show next steps
		fmt.Println(ui.SubtitleStyle.Render("Next steps:"))
		fmt.Printf("  • View build: amp builds get %s --agent %s\n", build.Name, agent)
		fmt.Printf("  • Monitor logs: amp builds logs %s --agent %s\n", build.Name, agent)
		fmt.Println()

		return nil
	},
}

var buildsLogsCmd = &cobra.Command{
	Use:   "logs <build-name>",
	Short: "Get build logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		buildName := args[0]

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agent, _ := cmd.Flags().GetString("agent")
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
		if agent == "" {
			return fmt.Errorf("agent name is required. Use --agent flag")
		}

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch logs from API
		logs, err := client.GetBuildLogs(org, project, agent, buildName)
		if err != nil {
			return fmt.Errorf("failed to get build logs: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(logs)
		}

		// Display logs
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Logs for Build: %s", ui.IconBuild, buildName)))
		fmt.Println()

		if len(logs.Logs) == 0 {
			fmt.Println(ui.RenderWarning("No logs available yet."))
			return nil
		}

		// Print each log entry with timestamp and level
		for _, entry := range logs.Logs {
			// Parse and format timestamp (expected RFC3339 format)
			timestamp := entry.Timestamp
			if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
				// Format as HH:MM:SS
				timestamp = t.Format("15:04:05")
			} else if strings.Contains(timestamp, "T") {
				// Fallback: best-effort extraction for non-RFC3339 timestamps
				parts := strings.SplitN(timestamp, "T", 2)
				if len(parts) == 2 {
					timePart := parts[1]
					// Strip trailing 'Z' if present
					if idx := strings.Index(timePart, "Z"); idx != -1 {
						timePart = timePart[:idx]
					}
					// Strip fractional seconds if present
					if idx := strings.Index(timePart, "."); idx != -1 {
						timePart = timePart[:idx]
					}
					timestamp = timePart
				}
			}

			levelPrefix := formatLogLevel(entry.LogLevel)
			if levelPrefix != "" {
				fmt.Printf("[%s] %s %s\n", timestamp, levelPrefix, entry.Log)
			} else {
				fmt.Printf("[%s] %s\n", timestamp, entry.Log)
			}
		}

		fmt.Println()
		return nil
	},
}

// Helper functions

// printBuildRow prints a styled key-value row for build details
func printBuildRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

// formatDuration calculates and formats duration between start and end times
func formatDuration(start time.Time, end *time.Time) string {
	var duration time.Duration
	if end == nil {
		// Build in progress - calculate from now
		duration = time.Since(start)
	} else {
		// Build completed - calculate actual duration
		duration = end.Sub(start)
	}
	return duration.Round(time.Second).String()
}

// truncateCommit shortens commit ID for display (first 8 chars)
func truncateCommit(commit string) string {
	if len(commit) > 8 {
		return commit[:8]
	}
	return commit
}

// formatLogLevel returns a styled prefix for log level (case-insensitive)
func formatLogLevel(level string) string {
	switch strings.ToLower(level) {
	case "error":
		return ui.ErrorStyle.Render("[ERROR]")
	case "warning":
		return ui.WarningStyle.Render("[WARN]")
	case "info":
		return ui.InfoStyle.Render("[INFO]")
	default:
		return ""
	}
}

func init() {
	// Register builds command with root
	rootCmd.AddCommand(buildsCmd)

	// Register subcommands with builds
	buildsCmd.AddCommand(buildsListCmd)
	buildsCmd.AddCommand(buildsGetCmd)
	buildsCmd.AddCommand(buildsTriggerCmd)
	buildsCmd.AddCommand(buildsLogsCmd)

	// Add --agent flag to all subcommands (required)
	buildsListCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	buildsGetCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	buildsTriggerCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	buildsLogsCmd.Flags().StringP("agent", "a", "", "Agent name (required)")

	// Add --commit flag to trigger command (optional)
	buildsTriggerCmd.Flags().StringP("commit", "c", "", "Commit ID (defaults to latest)")
}
