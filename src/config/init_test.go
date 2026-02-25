package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestInitTemplateConfigVariables(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `env:
  - name: CONFIG_PATH
    value: '{{ .ConfigPath }}'
  - name: CONFIG_DIR
    value: '{{ .ConfigDir }}'
  - name: IS_WSL
    value: '{{ .WSL }}'
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	escapedDir := strings.ReplaceAll(resolveConfigDir(configFile), `\`, `\\`)
	expected := "export CONFIG_PATH=\"" + configFile + "\"\n" +
		"export CONFIG_DIR=\"" + escapedDir + "\"\n" +
		fmt.Sprintf("export IS_WSL=\"%t\"", context.Current.WSL)
	assert.Equal(t, expected, script)
}

func TestInitAutoProgressWithWeights(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `progress:
  start_percentage: 10
  end_percentage: 70
alias:
  - name: ll
    value: ls -la
env:
  - name: AUTO_PROGRESS
    value: enabled
path:
  - value: /tmp/bin
script:
  - value: echo done
    weight: 2
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	got := Init(configFile, shell.BASH, true)

	expectedProgress := []string{
		`printf '\033]9;4;1;10\007'`,
		`printf '\033]9;4;1;22\007'`,
		`printf '\033]9;4;1;34\007'`,
		`printf '\033]9;4;1;46\007'`,
		`printf '\033]9;4;1;70\007'`,
	}

	currentIndex := 0
	for _, value := range expectedProgress {
		index := strings.Index(got[currentIndex:], value)
		assert.GreaterOrEqual(t, index, 0, "missing progress command %s", value)
		currentIndex += index + len(value)
	}
}
