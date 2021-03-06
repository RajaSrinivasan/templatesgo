package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/projtemplates/go/server/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Report the version of the application",
	Long: `
	Report the version of the application
	`,
	Run: Version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Version(cmd *cobra.Command, args []string) {
	version.Report()
}
