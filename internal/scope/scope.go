package scope

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/drakeafk/cmdsetgo/internal/events"
)

// GetGitRepoRoot returns the absolute path to the current Git repository root.
// If not in a Git repository, it returns an empty string and nil error.
func GetGitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			_ = exitErr
			return "", nil
		}
		return "", fmt.Errorf("failed to run git rev-parse: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// FilterEventsByRepoScope filters a slice of CmdEvent, returning only those
// whose current working directory (Cwd) is within the given repoRoot.
// If repoRoot is empty, it returns all events (global scope behavior).
func FilterEventsByRepoScope(evs []events.CmdEvent, repoRoot string) []events.CmdEvent {
	if repoRoot == "" {
		return evs // Global scope, no filtering
	}

	var filteredEvents []events.CmdEvent
	for _, event := range evs {
		if strings.HasPrefix(event.Cwd, repoRoot) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}

// FormatCwd returns a formatted string for the event's Cwd.
// If repoRoot is provided and the event's Cwd is within it,
// the path is made relative to the repoRoot. Otherwise, it returns
// the basename of the Cwd or the Cwd itself if it's very short.
func FormatCwd(cwd, repoRoot string) string {
	if repoRoot != "" && strings.HasPrefix(cwd, repoRoot) {
		relPath, err := filepath.Rel(repoRoot, cwd)
		if err == nil && relPath != "." { // If relPath is ".", it means cwd == repoRoot
			return relPath + "/"
		}
		return "repo/" // Default for repo root itself
	}
	// Global scope or not within repo, show last two segments or basename
	parts := strings.Split(cwd, string(filepath.Separator))
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], string(filepath.Separator)) + "/"
	}
	if len(parts) > 0 {
		return parts[len(parts)-1] + "/"
	}
	return cwd + "/" // Fallback for very short or root paths
}
