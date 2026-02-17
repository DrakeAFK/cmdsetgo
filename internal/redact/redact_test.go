package redact

import "testing"

func TestRedact(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		patterns []string
		want     string
	}{
		{
			"simple env var",
			"export GITHUB_TOKEN=secret123",
			nil,
			"export GITHUB_TOKEN=***REDACTED***",
		},
		{
			"quoted env var",
			"PASS=\"my password\" ls",
			nil,
			"PASS=\"***REDACTED***\" ls",
		},
		{
			"cli flag space",
			"mycli --token mytoken --command do-something",
			nil,
			"mycli --token ***REDACTED*** --command do-something",
		},
		{
			"cli flag equals",
			"mycli --password=mypassword",
			nil,
			"mycli --password=***REDACTED***",
		},
		{
			"custom regex",
			"echo 'sensitive: secret'",
			[]string{`secret`},
			"echo 'sensitive: ***REDACTED***'",
		},
		{
			"multiple secrets",
			"AWS_SECRET_ACCESS_KEY=abc GITHUB_TOKEN=xyz ./deploy",
			nil,
			"AWS_SECRET_ACCESS_KEY=***REDACTED*** GITHUB_TOKEN=***REDACTED*** ./deploy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Redact(tt.cmd, tt.patterns); got != tt.want {
				t.Errorf("Redact() = %v, want %v", got, tt.want)
			}
		})
	}
}
