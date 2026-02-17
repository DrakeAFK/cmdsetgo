package store

import (
	"os"
	"path/filepath"
)

const (
	DefaultEventsFile  = "events.jsonl"
	DefaultStateDir    = "state"
)

// GetConfigDir returns the default configuration directory ~/.cmdsetgo
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cmdsetgo"), nil
}

// GetEventsPath returns the path to the events file.
// If the CMDSETGO_EVENTS_PATH env var is set, it uses that.
// Otherwise, it defaults to ~/.cmdsetgo/events.jsonl
func GetEventsPath() (string, error) {
	if path := os.Getenv("CMDSETGO_EVENTS_PATH"); path != "" {
		return path, nil
	}
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, DefaultEventsFile), nil
}

// GetStateDir returns the path to the state directory ~/.cmdsetgo/state/
func GetStateDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, DefaultStateDir), nil
}

// EnsureDirs creates the config and state directories if they don't exist.
func EnsureDirs() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	
	stateDir, err := GetStateDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.MkdirAll(stateDir, 0755)
}
