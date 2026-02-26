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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip check for install, status, and help commands
		if cmd.Name() == "install" || cmd.Name() == "status" || cmd.Name() == "help" || cmd.Name() == "cmdsetgo" || cmd.Name() == "uninstall" {
			return
		}

		if os.Getenv("CMDSETGO_EVENTS_PATH") == "" {
			fmt.Println("Note: cmdsetgo hook is not active in this session.")
			fmt.Println("Run `cmdsetgo install` to set it up, or restart your terminal if you just installed it.")
			fmt.Println()
		}
	},
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
