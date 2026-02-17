package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdsetgo",
	Short: "cmdsetgo - Turn terminal chaos into a clean script",
	Long: `cmdsetgo records terminal commands into a structured log and lets you 
view, pick, and export them as a clean bash script or markdown runbook.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags can be added here
}
