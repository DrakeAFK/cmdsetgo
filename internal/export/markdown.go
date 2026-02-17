package export

import (
	"fmt"
	"io"
	"time"

	"github.com/drakeafk/cmdsetgo/internal/pick"
	"github.com/drakeafk/cmdsetgo/internal/redact"
)

// MarkdownExporter generates a markdown runbook.
func MarkdownExporter(w io.Writer, selection pick.Selection, redactRegex []string) error {
	fmt.Fprintln(w, "# cmdsetgo runbook")
	fmt.Fprintf(w, "\nGenerated at %s  \n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(w, "Scope: `%s`  \n", selection.Scope)
	if selection.RepoRoot != "" {
		fmt.Fprintf(w, "Repo Root: `%s`  \n", selection.RepoRoot)
	}
	fmt.Fprintln(w)

	currentCwd := ""
	for _, ev := range selection.Items {
		if ev.Cwd != currentCwd {
			fmt.Fprintf(w, "## In `%s`\n\n", ev.Cwd)
			currentCwd = ev.Cwd
		}

		cmdRedacted := redact.Redact(ev.Cmd, redactRegex)
		fmt.Fprintln(w, "```bash")
		fmt.Fprintf(w, "# %s\n", ev.Ts.Format(time.RFC3339))
		fmt.Fprintln(w, cmdRedacted)
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w)
	}

	return nil
}
