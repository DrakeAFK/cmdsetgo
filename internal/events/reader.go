package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ReadEvents reads CmdEvent objects from a JSONL file.
// It returns a slice of CmdEvent and an error if the file cannot be opened.
// Malformed lines are skipped with a warning.
func ReadEvents(filePath string) ([]CmdEvent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open events file: %w", err)
	}
	defer file.Close()

	var events []CmdEvent
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var event CmdEvent
		if err := json.Unmarshal(line, &event); err != nil {
			log.Printf("Warning: Skipping malformed event line in %s: %v, line: %s", filePath, err, string(line))
			continue
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading events file: %w", err)
	}

	return events, nil
}
