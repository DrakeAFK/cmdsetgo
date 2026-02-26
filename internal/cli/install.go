package cli

import (
	"fmt"
	"os"

	"github.com/drakeafk/cmdsetgo/internal/shell"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var (
	installShell  string
	installEvents string
	installBin    string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the shell hook to record commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		shellName := installShell
		if shellName == "" {
			shellName = shell.DetectShell()
			if shellName == "" {
				return fmt.Errorf("could not auto-detect shell; please specify with --shell bash|zsh")
			}
			fmt.Printf("Detected shell: %s\n", shellName)
		}

		eventsPath := installEvents
		if eventsPath == "" {
			var err error
			eventsPath, err = store.GetEventsPath()
			if err != nil {
				return err
			}
		}

		// Ensure config dir exists
		if err := store.EnsureDirs(); err != nil {
			return err
		}

		binaryPath := installBin
		if binaryPath == "" {
			// Try to get absolute path to current executable
			exe, err := os.Executable()
			if err == nil {
				binaryPath = exe
			}
		}

		if err := shell.Install(shellName, eventsPath, binaryPath); err != nil {
			return err
		}

		fmt.Printf("Successfully installed cmdsetgo hook for %s.\n", shellName)
		if binaryPath != "" {
			fmt.Printf("Added alias: cmdsetgo -> %s\n", binaryPath)
		}
		fmt.Println("Please restart your terminal or source your rc file to start recording.")
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the shell hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if installShell == "" {
			return fmt.Errorf("shell is required (--shell bash|zsh)")
		}

		if err := shell.Uninstall(installShell); err != nil {
			return err
		}

		fmt.Printf("Successfully uninstalled cmdsetgo hook for %s.\n", installShell)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)

	installCmd.Flags().StringVar(&installShell, "shell", "", "Shell to install hook for (bash or zsh)")
	installCmd.Flags().StringVar(&installEvents, "events", "", "Path to the events log file (optional)")
	installCmd.Flags().StringVar(&installBin, "bin", "", "Absolute path to the cmdsetgo binary (optional)")

	uninstallCmd.Flags().StringVar(&installShell, "shell", "", "Shell to uninstall hook from (bash or zsh)")
}
