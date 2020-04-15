package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details ",
	Long: `
	Shows status of the server operation.`,
	Run: Show,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func Show(cmd *cobra.Command, args []string) {
	fmt.Println("show called")
}
