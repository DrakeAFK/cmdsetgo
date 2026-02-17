package events

import (
	"encoding/json"
	"os"
)

// WriteEvent appends a single CmdEvent to the specified JSONL file.
func WriteEvent(filePath string, event CmdEvent) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(event)
}
