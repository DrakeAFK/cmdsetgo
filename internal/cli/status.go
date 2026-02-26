package cli

import (
	"fmt"
	"os"

	"github.com/drakeafk/cmdsetgo/internal/shell"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the installation status and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cmdsetgo status:")

		// 1. Check if hook is active in current session
		eventsPath := os.Getenv("CMDSETGO_EVENTS_PATH")
		if eventsPath != "" {
			fmt.Printf("  - Active in current session: YES\n")
			fmt.Printf("  - Events log location: %s\n", eventsPath)
		} else {
			fmt.Printf("  - Active in current session: NO (Hook not detected in this terminal)\n")
		}

		// 2. Check if installed in RC files
		detectedShell := shell.DetectShell()
		isInstalled := false
		if detectedShell != "" {
			fmt.Printf("  - Detected shell: %s\n", detectedShell)
			var err error
			isInstalled, err = shell.IsInstalled(detectedShell)
			if err != nil {
				fmt.Printf("  - Installed in %s RC: Error checking (%v)\n", detectedShell, err)
			} else if isInstalled {
				fmt.Printf("  - Installed in %s RC: YES\n", detectedShell)
			} else {
				fmt.Printf("  - Installed in %s RC: NO (Run `cmdsetgo install` to set up)\n", detectedShell)
			}
		} else {
			fmt.Printf("  - Detected shell: Could not detect\n")
		}

		// 3. Detect stale session (active in memory but uninstalled)
		if eventsPath != "" && !isInstalled {
			fmt.Println()
			fmt.Println("⚠️  STALE SESSION DETECTED:")
			fmt.Println("   The cmdsetgo hook is active in this terminal, but it has been uninstalled from your configuration.")
			fmt.Println("   Please close this terminal tab and open a new one to fully deactivate it.")
		}

		// 3. Storage location
		baseDir, _ := store.GetConfigDir()
		fmt.Printf("  - Config/Data directory: %s\n", baseDir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
