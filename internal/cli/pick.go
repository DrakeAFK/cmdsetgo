package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/drakeafk/cmdsetgo/internal/events"
	"github.com/drakeafk/cmdsetgo/internal/pick"
	"github.com/drakeafk/cmdsetgo/internal/scope"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var (
	pickNum       int
	pickScope     string
	excludeCommon bool
	excludeRegex  []string
)

var pickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Interactively pick and reorder commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		eventsPath, err := store.GetEventsPath()
		if err != nil {
			return err
		}

		allEvents, err := events.ReadEvents(eventsPath)
		if err != nil {
			return err
		}

		var repoRoot string
		if pickScope == "repo" || (pickScope == "" && isInGitRepo()) {
			repoRoot, err = scope.GetGitRepoRoot()
			if err != nil {
				return err
			}
		}

		filtered := scope.FilterEventsByRepoScope(allEvents, repoRoot)

		patterns := excludeRegex
		if excludeCommon {
			patterns = append(patterns, pick.CommonExclusions...)
		}
		filtered = pick.FilterExclusions(filtered, patterns)

		// Take last N
		if len(filtered) > pickNum {
			filtered = filtered[len(filtered)-pickNum:]
		}

		if len(filtered) == 0 {
			fmt.Println("No commands found in this scope.")
			return nil
		}

		// Print list
		printTable(filtered, repoRoot)

		fmt.Print("\nSelect commands in the order you want (e.g. \"5 2 3\", \"1-4 7\", or \"all\"): ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return nil
		}
		input := scanner.Text()

		indices, err := pick.ParseSelection(input, len(filtered))
		if err != nil {
			return err
		}

		if len(indices) == 0 {
			fmt.Println("No commands selected.")
			return nil
		}

		var selectedItems []events.CmdEvent
		for _, idx := range indices {
			selectedItems = append(selectedItems, filtered[idx-1])
		}

		// Save selection
		selectionID := time.Now().Format("20060102-150405")
		selection := pick.Selection{
			ID:        selectionID,
			CreatedAt: time.Now().Format(time.RFC3339),
			Scope:     pickScope,
			RepoRoot:  repoRoot,
			Items:     selectedItems,
		}

		stateDir, err := store.GetStateDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(stateDir, 0755); err != nil {
			return err
		}

		selectionPath := filepath.Join(stateDir, fmt.Sprintf("selection-%s.json", selectionID))
		file, err := os.Create(selectionPath)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(selection); err != nil {
			return err
		}

		fmt.Printf("\nSaved selection: %s\n", selectionID)
		fmt.Printf("Export: cmdsetgo export --selection %s --format bash --out run.sh\n", selectionID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pickCmd)
	pickCmd.Flags().IntVarP(&pickNum, "num", "n", 50, "Number of commands to show")
	pickCmd.Flags().StringVar(&pickScope, "scope", "", "Scope: repo or global (default auto-detect)")
	pickCmd.Flags().BoolVar(&excludeCommon, "exclude-common", true, "Exclude common noise commands like ls, cd, etc.")
	pickCmd.Flags().StringSliceVar(&excludeRegex, "exclude-regex", []string{}, "Regex patterns to exclude commands")
}
