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
	"github.com/Kavirubc/wso2-amp-cli/internal/util"
	"github.com/spf13/cobra"
)

var tracesCmd = &cobra.Command{
	Use:   "traces",
	Short: "View distributed traces for agents",
	Long:  `Commands for listing, viewing, and exporting distributed traces from AI agents.`,
}

var tracesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List traces for an agent",
	Long: `List distributed traces for a deployed agent.

Examples:
  amp traces list --agent myagent --env development
  amp traces list --agent myagent --env dev --since 24h
  amp traces list --agent myagent --env dev --limit 50
  amp traces list --agent myagent --env dev --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
		since, _ := cmd.Flags().GetString("since")
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

		// Build trace list options
		opts := api.TraceListOptions{
			Environment: envName,
			Limit:       limit,
			SortOrder:   sortOrder,
		}

		// Parse --since flag into start time (default to 1h if not provided)
		sinceValue := since
		if sinceValue == "" {
			sinceValue = "1h"
		}
		startTime, err := util.ParseSinceDuration(sinceValue)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		opts.StartTime = startTime.Format(time.RFC3339)
		opts.EndTime = time.Now().Format(time.RFC3339)

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch traces from API
		traces, err := client.ListTraces(org, project, agentName, opts)
		if err != nil {
			return fmt.Errorf("failed to list traces: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(traces)
		}

		// Display traces
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Traces: %s (%s)", ui.IconTrace, agentName, envName)))
		fmt.Println()

		if len(traces.Traces) == 0 {
			fmt.Println(ui.RenderWarning("No traces found for the specified criteria."))
			return nil
		}

		// Build table data
		headers := []string{"TRACE ID", "ROOT SPAN", "STATUS", "DURATION", "TIMESTAMP"}
		rows := make([][]string, len(traces.Traces))
		for i, trace := range traces.Traces {
			errorCount := 0
			if trace.Status != nil {
				errorCount = trace.Status.ErrorCount
			}
			status := ui.TraceStatusCell(errorCount)

			rows[i] = []string{
				ui.TruncateTraceID(trace.TraceID),
				ui.TruncateString(trace.RootSpanName, 30),
				status,
				ui.FormatNanosDuration(trace.DurationInNanos),
				ui.FormatTraceTimestamp(trace.StartTime),
			}
		}

		// Render styled table
		fmt.Println(ui.RenderTable(headers, rows))
		fmt.Println()
		fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("Showing %d of %d traces", len(traces.Traces), traces.TotalCount)))

		return nil
	},
}

var tracesGetCmd = &cobra.Command{
	Use:   "get <traceId>",
	Short: "Get details of a specific trace",
	Long: `Get detailed information about a specific trace including its span tree.

Examples:
  amp traces get abc123def456 --agent myagent --env development
  amp traces get abc123def456 --agent myagent --env dev --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		traceID := args[0]

		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
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

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Fetch trace details from API
		traceDetails, err := client.GetTrace(org, project, agentName, traceID, envName)
		if err != nil {
			return fmt.Errorf("failed to get trace: %w", err)
		}

		// JSON output
		if output == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(traceDetails)
		}

		// Display trace details
		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("%s Trace: %s", ui.IconTrace, ui.TruncateTraceID(traceID))))
		fmt.Println()

		if len(traceDetails.Spans) == 0 {
			fmt.Println(ui.RenderWarning("No spans found for this trace."))
			return nil
		}

		// Show trace summary
		printTraceRow("Trace ID:", traceID)
		printTraceRow("Environment:", envName)
		printTraceRow("Total Spans:", fmt.Sprintf("%d", traceDetails.TotalCount))
		fmt.Println()

		// Build and display span tree
		fmt.Println(ui.SubtitleStyle.Render("  Span Tree:"))
		fmt.Println()
		roots := ui.BuildSpanTree(traceDetails.Spans)
		count := 0
		maxSpans := 50 // Truncate at 50 spans
		treeOutput := ui.RenderSpanTree(roots, "  ", maxSpans, &count)
		fmt.Print(treeOutput)
		fmt.Println()

		return nil
	},
}

var tracesExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export traces with full span details",
	Long: `Export traces with complete span information in JSON format.

Examples:
  amp traces export --agent myagent --env development --since 24h
  amp traces export --agent myagent --env dev --since 7d --file traces.json
  amp traces export --agent myagent --env dev --limit 200 --file traces.json --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		org, _ := cmd.Flags().GetString("org")
		project, _ := cmd.Flags().GetString("project")
		agentName, _ := cmd.Flags().GetString("agent")
		envName, _ := cmd.Flags().GetString("env")
		since, _ := cmd.Flags().GetString("since")
		limit, _ := cmd.Flags().GetInt("limit")
		filePath, _ := cmd.Flags().GetString("file")
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

		// Check if file exists and --force not set
		if filePath != "" && !force {
			if _, err := os.Stat(filePath); err == nil {
				return fmt.Errorf("file %q already exists. Use --force to overwrite", filePath)
			}
		}

		// Build trace list options
		opts := api.TraceListOptions{
			Environment: envName,
			Limit:       limit,
		}

		// Parse --since flag into start time (default to 24h if not provided)
		sinceValue := since
		if sinceValue == "" {
			sinceValue = "24h"
		}
		startTime, err := util.ParseSinceDuration(sinceValue)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		opts.StartTime = startTime.Format(time.RFC3339)
		opts.EndTime = time.Now().Format(time.RFC3339)

		// Create API client
		client := api.NewClient(
			config.GetAPIURL(),
			config.GetAPIKeyHeader(),
			config.GetAPIKeyValue(),
		)

		// Show progress
		fmt.Printf("Exporting traces from %s (%s)...\n", agentName, envName)

		// Fetch traces from API
		traces, err := client.ExportTraces(org, project, agentName, opts)
		if err != nil {
			return fmt.Errorf("failed to export traces: %w", err)
		}

		if len(traces.Traces) == 0 {
			fmt.Println(ui.RenderWarning("No traces found for the specified criteria."))
			return nil
		}

		// Encode to JSON
		var jsonData []byte
		jsonData, err = json.MarshalIndent(traces, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode traces: %w", err)
		}

		// Write to file or stdout
		if filePath != "" {
			err = os.WriteFile(filePath, jsonData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("Exported %d traces to %s", len(traces.Traces), filePath)))
		} else {
			fmt.Println(string(jsonData))
		}

		return nil
	},
}

// printTraceRow prints a styled key-value row for trace details
func printTraceRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	rootCmd.AddCommand(tracesCmd)
	tracesCmd.AddCommand(tracesListCmd)
	tracesCmd.AddCommand(tracesGetCmd)
	tracesCmd.AddCommand(tracesExportCmd)

	// Add flags for list command
	tracesListCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	tracesListCmd.Flags().StringP("env", "e", "", "Environment name (required)")
	tracesListCmd.Flags().String("since", "", "Show traces since duration (e.g., 1h, 24h, 7d). Default: 1h")
	tracesListCmd.Flags().Int("limit", 25, "Maximum number of traces to return (1-1000)")
	tracesListCmd.Flags().String("sort", "desc", "Sort order (asc/desc)")

	// Add flags for get command
	tracesGetCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	tracesGetCmd.Flags().StringP("env", "e", "", "Environment name (required)")

	// Add flags for export command
	tracesExportCmd.Flags().StringP("agent", "a", "", "Agent name (required)")
	tracesExportCmd.Flags().StringP("env", "e", "", "Environment name (required)")
	tracesExportCmd.Flags().String("since", "", "Export traces since duration (e.g., 1h, 24h, 7d). Default: 24h")
	tracesExportCmd.Flags().StringP("file", "f", "", "Output file path (outputs to stdout if not specified)")
	tracesExportCmd.Flags().Bool("force", false, "Overwrite existing file")
	tracesExportCmd.Flags().Int("limit", 100, "Maximum number of traces to export (1-1000)")
}
