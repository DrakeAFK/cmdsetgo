package redact

import (
	"regexp"
	"strings"
)

var SecretKeys = []string{
	"TOKEN", "KEY", "SECRET", "PASSWORD", "PASS", "AUTH", "BEARER",
	"AWS_SECRET_ACCESS_KEY", "GITHUB_TOKEN",
}

var SecretFlags = []string{
	"token", "password", "apikey", "api-key",
}

// Redact replaces sensitive values in command strings with ***REDACTED***.
func Redact(cmd string, customRegexes []string) string {
	redacted := cmd

	// 1. Redact KEY=VALUE
	for _, key := range SecretKeys {
		// Match KEY=VALUE where VALUE is not already redacted.
		// Handles shell exports and inline env vars.
		re := regexp.MustCompile(identRegex(key) + `=[^ \t\n\r"']+`)
		redacted = re.ReplaceAllStringFunc(redacted, func(match string) string {
			parts := strings.SplitN(match, "=", 2)
			return parts[0] + "=***REDACTED***"
		})

		// Handle quoted values KEY="VALUE"
		reQuoted := regexp.MustCompile(identRegex(key) + `="[^"]+"`)
		redacted = reQuoted.ReplaceAllStringFunc(redacted, func(match string) string {
			parts := strings.SplitN(match, "=", 2)
			return parts[0] + "=\"***REDACTED***\""
		})
	}

	// 2. Redact CLI flags --token <value>
	for _, flag := range SecretFlags {
		re := regexp.MustCompile(`--` + flag + `[= ]+[^ \t\n\r"']+`)
		redacted = re.ReplaceAllStringFunc(redacted, func(match string) string {
			// Find separator (space or =)
			sep := " "
			if strings.Contains(match, "=") {
				sep = "="
			}
			parts := strings.SplitN(match, sep, 2)
			return parts[0] + sep + "***REDACTED***"
		})
	}

	// 3. Apply custom regexes
	for _, pattern := range customRegexes {
		if re, err := regexp.Compile(pattern); err == nil {
			redacted = re.ReplaceAllString(redacted, "***REDACTED***")
		}
	}

	return redacted
}

func identRegex(key string) string {
	// Case insensitive match for the key, but must be exactly the key or end with it (e.g. GITHUB_TOKEN)
	return `(?i:[A-Z0-9_]*` + key + `)`
}
