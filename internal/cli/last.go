package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/drakeafk/cmdsetgo/internal/events"
	"github.com/drakeafk/cmdsetgo/internal/scope"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var (
	lastNum    int
	lastScope  string
	lastFormat string
)

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "View the last N commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		eventsPath, err := store.GetEventsPath()
		if err != nil {
			return err
		}

		allEvents, err := events.ReadEvents(eventsPath)
		if err != nil {
			// If file doesn't exist, just treat as empty
			if !os.IsNotExist(err) {
				return err
			}
			allEvents = []events.CmdEvent{}
		}

		var repoRoot string
		if lastScope == "repo" || (lastScope == "" && isInGitRepo()) {
			repoRoot, err = scope.GetGitRepoRoot()
			if err != nil {
				return err
			}
		}

		filtered := scope.FilterEventsByRepoScope(allEvents, repoRoot)

		// Take last N
		if len(filtered) > lastNum {
			filtered = filtered[len(filtered)-lastNum:]
		}

		if lastFormat == "json" {
			return printJSON(filtered)
		}

		printTable(filtered, repoRoot)
		return nil
	},
}

func isInGitRepo() bool {
	root, err := scope.GetGitRepoRoot()
	return err == nil && root != ""
}

func printJSON(evs []events.CmdEvent) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(evs)
}

func printTable(evs []events.CmdEvent, repoRoot string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for i, ev := range evs {
		formattedTime := ev.Ts.Local().Format("15:04:05")
		shortCwd := scope.FormatCwd(ev.Cwd, repoRoot)
		fmt.Fprintf(w, "# %d\t%s\t%s\t%s\t(%d)\n", i+1, formattedTime, shortCwd, ev.Cmd, ev.Exit)
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.Flags().IntVarP(&lastNum, "num", "n", 30, "Number of commands to show")
	lastCmd.Flags().StringVar(&lastScope, "scope", "", "Scope: repo or global (default auto-detect)")
	lastCmd.Flags().StringVar(&lastFormat, "format", "table", "Output format: table or json")
}
