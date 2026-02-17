package scope

import "testing"

func TestFormatCwd(t *testing.T) {
	tests := []struct {
		cwd      string
		repoRoot string
		want     string
	}{
		{"/Users/me/code/proj/src", "/Users/me/code/proj", "src/"},
		{"/Users/me/code/proj", "/Users/me/code/proj", "repo/"},
		{"/Users/me/code/other", "/Users/me/code/proj", "code/other/"},
		{"/tmp", "", "tmp/"},
		{"/", "", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.cwd, func(t *testing.T) {
			if got := FormatCwd(tt.cwd, tt.repoRoot); got != tt.want {
				t.Errorf("FormatCwd() = %v, want %v", got, tt.want)
			}
		})
	}
}
