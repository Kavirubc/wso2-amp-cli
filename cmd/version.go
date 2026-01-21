package cmd

import (
	"fmt"
	"runtime"

	"github.com/Kavirubc/wso2-amp-cli/internal/ui"
	"github.com/spf13/cobra"
)

// Version information set via ldflags at build time
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, build date, git commit, Go version, and platform information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.TitleStyle.Render("amp CLI"))
		fmt.Println()
		printVersionRow("Version:", Version)
		printVersionRow("Git Commit:", GitCommit)
		printVersionRow("Build Date:", BuildDate)
		printVersionRow("Go Version:", runtime.Version())
		printVersionRow("Platform:", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		fmt.Println()
	},
}

// printVersionRow prints a styled key-value row for version info
func printVersionRow(key, value string) {
	fmt.Printf("  %s  %s\n", ui.KeyStyle.Render(key), ui.ValueStyle.Render(value))
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
