package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestShouldSkipInitOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ttyOnly  bool
		stdinTTY bool
		wantSkip bool
	}{
		{
			name:     "tty-only disabled",
			ttyOnly:  false,
			stdinTTY: false,
			wantSkip: false,
		},
		{
			name:     "interactive terminal",
			ttyOnly:  true,
			stdinTTY: true,
			wantSkip: false,
		},
		{
			name:     "piped output but interactive input should still print",
			ttyOnly:  true,
			stdinTTY: true,
			wantSkip: false,
		},
		{
			name:     "non-interactive process",
			ttyOnly:  true,
			stdinTTY: false,
			wantSkip: true,
		},
		{
			name:     "piped input should skip even when output is interactive",
			ttyOnly:  true,
			stdinTTY: false,
			wantSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldSkipInitOutput(tt.ttyOnly, tt.stdinTTY)
			if got != tt.wantSkip {
				t.Fatalf("shouldSkipInitOutput() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}

func TestRunInitWritesToStdout(t *testing.T) {
	originalConfig := config
	originalPrintOutput := printOutput
	originalTTYOnly := ttyOnly
	t.Cleanup(func() {
		config = originalConfig
		printOutput = originalPrintOutput
		ttyOnly = originalTTYOnly
	})

	configFile := filepath.Join(t.TempDir(), "aliae.yaml")
	err := os.WriteFile(configFile, []byte(`alias:
  - name: ll
    value: ls -la
`), 0o600)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	config = configFile
	printOutput = true
	ttyOnly = false

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	runInit(cmd, "zsh")

	assert.NotEmpty(t, strings.TrimSpace(stdout.String()))
	assert.Empty(t, stderr.String())
}
