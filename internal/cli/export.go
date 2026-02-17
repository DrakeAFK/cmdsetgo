package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/drakeafk/cmdsetgo/internal/export"
	"github.com/drakeafk/cmdsetgo/internal/pick"
	"github.com/drakeafk/cmdsetgo/internal/store"
	"github.com/spf13/cobra"
)

var (
	exportFormat    string
	exportOut       string
	exportSelection string
	exportRedact    []string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export selected commands to a script or runbook",
	RunE: func(cmd *cobra.Command, args []string) error {
		stateDir, err := store.GetStateDir()
		if err != nil {
			return err
		}

		selectionPath := exportSelection
		if selectionPath == "" {
			// Find most recent selection
			selectionPath, err = findMostRecentSelection(stateDir)
			if err != nil {
				return err
			}
		} else if !filepath.IsAbs(selectionPath) && !strings.Contains(selectionPath, string(filepath.Separator)) {
			// If it's just an ID, assume it's in stateDir
			selectionPath = filepath.Join(stateDir, fmt.Sprintf("selection-%s.json", selectionPath))
		}

		file, err := os.Open(selectionPath)
		if err != nil {
			return fmt.Errorf("failed to open selection file %s: %w", selectionPath, err)
		}
		defer file.Close()

		var selection pick.Selection
		if err := json.NewDecoder(file).Decode(&selection); err != nil {
			return fmt.Errorf("failed to decode selection file: %w", err)
		}

		var out io.Writer = os.Stdout
		if exportOut != "" {
			f, err := os.Create(exportOut)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}

		switch exportFormat {
		case "bash":
			return export.BashExporter(out, selection, exportRedact)
		case "md", "markdown":
			return export.MarkdownExporter(out, selection, exportRedact)
		default:
			return fmt.Errorf("unknown format: %s", exportFormat)
		}
	},
}

func findMostRecentSelection(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var selections []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "selection-") && strings.HasSuffix(f.Name(), ".json") {
			selections = append(selections, f.Name())
		}
	}

	if len(selections) == 0 {
		return "", fmt.Errorf("no selections found in %s", dir)
	}

	// Sort descending (most recent first)
	sort.Slice(selections, func(i, j int) bool {
		return selections[i] > selections[j]
	})

	return filepath.Join(dir, selections[0]), nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportFormat, "format", "bash", "Output format: bash or md")
	exportCmd.Flags().StringVar(&exportOut, "out", "", "Output file path (default stdout)")
	exportCmd.Flags().StringVar(&exportSelection, "selection", "", "Selection ID or path to selection file")
	exportCmd.Flags().StringSliceVar(&exportRedact, "redact-regex", []string{}, "Custom regex patterns to redact")
}
