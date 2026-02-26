package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	StartMarker = "# >>> cmdsetgo >>>"
	EndMarker   = "# <<< cmdsetgo <<<"
)

const BashHook = `
# cmdsetgo bash hook
_cmdsetgo_hook() {
    local exit_code=$?
    local last_cmd=$(history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
    
    # Avoid logging cmdsetgo commands themselves to keep things clean
    if [[ "$last_cmd" == cmdsetgo* ]]; then
        return
    fi

    local ts=$(date +"%Y-%m-%dT%H:%M:%S%z")
    local cwd="$PWD"
    local user="$USER"
    local host="$HOSTNAME"
    
    printf '{"type":"cmd","ts":"%s","shell":"bash","host":"%s","user":"%s","cwd":"%s","cmd":"%s","exit":%d}\n' \
        "$ts" "$host" "$user" "$cwd" "$(echo "$last_cmd" | sed 's/"/\\"/g')" "$exit_code" >> "$CMDSETGO_EVENTS_PATH"
}
if [[ ! "$PROMPT_COMMAND" =~ _cmdsetgo_hook ]]; then
    PROMPT_COMMAND="_cmdsetgo_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi
`

const ZshHook = `
# cmdsetgo zsh hook
_cmdsetgo_preexec() {
    _CMDSETGO_LAST_CMD="$1"
}
_cmdsetgo_precmd() {
    local exit_code=$?
    if [[ -z "$_CMDSETGO_LAST_CMD" ]]; then
        return
    fi
    if [[ "$_CMDSETGO_LAST_CMD" == cmdsetgo* ]]; then
        unset _CMDSETGO_LAST_CMD
        return
    fi

    local ts=$(date +"%Y-%m-%dT%H:%M:%S%z")
    local cwd="$PWD"
    local user="$USER"
    local host="$HOST"
    
    printf '{"type":"cmd","ts":"%s","shell":"zsh","host":"%s","user":"%s","cwd":"%s","cmd":"%s","exit":%d}\n' \
        "$ts" "$host" "$user" "$cwd" "$(echo "$_CMDSETGO_LAST_CMD" | sed 's/"/\\"/g')" "$exit_code" >> "$CMDSETGO_EVENTS_PATH"
    unset _CMDSETGO_LAST_CMD
}
autoload -Uz add-zsh-hook
add-zsh-hook preexec _cmdsetgo_preexec
add-zsh-hook precmd _cmdsetgo_precmd
`

func GetRCPath(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch shell {
	case "bash":
		paths := []string{".bashrc", ".bash_profile"}
		for _, p := range paths {
			path := filepath.Join(home, p)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
		return filepath.Join(home, ".bashrc"), nil
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// DetectShell attempts to detect the current user shell.
func DetectShell() string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv != "" {
		if strings.Contains(shellEnv, "zsh") {
			return "zsh"
		}
		if strings.Contains(shellEnv, "bash") {
			return "bash"
		}
	}
	return ""
}

// IsInstalled checks if the cmdsetgo hook is installed in the shell's RC file.
func IsInstalled(shellName string) (bool, error) {
	rcPath, err := GetRCPath(shellName)
	if err != nil {
		return false, err
	}

	content, err := os.ReadFile(rcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return strings.Contains(string(content), StartMarker), nil
}

func Install(shellName string, eventsPath string, binaryPath string) error {
	rcPath, err := GetRCPath(shellName)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(rcPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	sContent := string(content)
	if strings.Contains(sContent, StartMarker) {
		// Already installed, remove old one first (idempotent)
		if err := Uninstall(shellName); err != nil {
			return err
		}
		// Refresh content
		content, _ = os.ReadFile(rcPath)
		sContent = string(content)
	}

	hook := BashHook
	if shellName == "zsh" {
		hook = ZshHook
	}

	aliasLine := ""
	if binaryPath != "" {
		aliasLine = fmt.Sprintf("alias cmdsetgo=\"%s\"\n", binaryPath)
	}

	block := fmt.Sprintf("\n%s\nexport CMDSETGO_EVENTS_PATH=\"%s\"\n%s%s\n%s\n", StartMarker, eventsPath, aliasLine, hook, EndMarker)

	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(block)
	return err
}

func Uninstall(shellName string) error {
	rcPath, err := GetRCPath(shellName)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(rcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	sContent := string(content)
	if !strings.Contains(sContent, StartMarker) {
		return nil
	}

	lines := strings.Split(sContent, "\n")
	var newLines []string
	inBlock := false
	for _, line := range lines {
		if strings.TrimSpace(line) == StartMarker {
			inBlock = true
			continue
		}
		if strings.TrimSpace(line) == EndMarker {
			inBlock = false
			continue
		}
		if !inBlock {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(rcPath, []byte(strings.Join(newLines, "\n")), 0644)
}
