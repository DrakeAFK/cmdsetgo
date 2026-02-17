package cli

import (
	"fmt"

	"github.com/drakeafk/cmdsetgo/internal/shell"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var (
	installShell  string
	installEvents string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the shell hook to record commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if installShell == "" {
			return fmt.Errorf("shell is required (--shell bash|zsh)")
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

		if err := shell.Install(installShell, eventsPath); err != nil {
			return err
		}

		fmt.Printf("Successfully installed cmdsetgo hook for %s.\n", installShell)
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

	uninstallCmd.Flags().StringVar(&installShell, "shell", "", "Shell to uninstall hook from (bash or zsh)")
}
