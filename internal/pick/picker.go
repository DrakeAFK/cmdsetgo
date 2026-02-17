package pick

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/drakeafk/cmdsetgo/internal/events"
)

var CommonExclusions = []string{
	`^(ls|cd|pwd|clear|exit)$`,
}

// FilterExclusions applies regex filters to the command string of each event.
func FilterExclusions(evs []events.CmdEvent, patterns []string) []events.CmdEvent {
	var filtered []events.CmdEvent
	regexes := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			regexes = append(regexes, re)
		}
	}

	for _, ev := range evs {
		excluded := false
		// Normalize command for simple common exclusions (check first word)
		cmdFirstWord := strings.Fields(ev.Cmd)
		var cmdBase string
		if len(cmdFirstWord) > 0 {
			cmdBase = cmdFirstWord[0]
		}

		for _, re := range regexes {
			if re.MatchString(ev.Cmd) || (cmdBase != "" && re.MatchString(cmdBase)) {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, ev)
		}
	}
	return filtered
}

// ParseSelection parses a string like "1 3-5 2" into a list of 1-based indices.
// It returns the list of indices in the order they were specified.
func ParseSelection(input string, maxIndex int) ([]int, error) {
	tokens := strings.Fields(strings.ReplaceAll(input, ",", " "))
	var selection []int
	seen := make(map[int]bool)

	for _, token := range tokens {
		if token == "all" {
			for i := 1; i <= maxIndex; i++ {
				if !seen[i] {
					selection = append(selection, i)
					seen[i] = true
				}
			}
			continue
		}

		if strings.Contains(token, "-") {
			parts := strings.Split(token, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", token)
			}
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range: %s", token)
			}
			if start > end {
				return nil, fmt.Errorf("start of range cannot be greater than end: %s", token)
			}
			for i := start; i <= end; i++ {
				if i < 1 || i > maxIndex {
					return nil, fmt.Errorf("index out of range: %d", i)
				}
				if !seen[i] {
					selection = append(selection, i)
					seen[i] = true
				}
			}
		} else {
			idx, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("invalid index: %s", token)
			}
			if idx < 1 || idx > maxIndex {
				return nil, fmt.Errorf("index out of range: %d", idx)
			}
			if !seen[idx] {
				selection = append(selection, idx)
				seen[idx] = true
			}
		}
	}

	return selection, nil
}

type Selection struct {
	ID        string            `json:"id"`
	CreatedAt string            `json:"created_at"`
	Scope     string            `json:"scope"`
	RepoRoot  string            `json:"repo_root"`
	Items     []events.CmdEvent `json:"items"`
}
