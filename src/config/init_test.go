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
  - name: DOTFILES_DIR
    value: '{{ .Env.DOTFILES }}'
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)
	t.Setenv("DOTFILES", "/tmp/dotfiles")

	script := Init(configFile, shell.BASH, true)
	escapedDir := strings.ReplaceAll(resolveConfigDir(configFile), `\`, `\\`)
	expected := "export CONFIG_PATH=\"" + configFile + "\"\n" +
		"export CONFIG_DIR=\"" + escapedDir + "\"\n" +
		fmt.Sprintf("export IS_WSL=\"%t\"\n", context.Current.WSL) +
		"export DOTFILES_DIR=\"/tmp/dotfiles\""
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

func TestInitAutoProgressWithFractionalScriptWeight(t *testing.T) {
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
    weight: 0.5
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	got := Init(configFile, shell.BASH, true)

	expectedProgress := []string{
		`printf '\033]9;4;1;10\007'`,
		`printf '\033]9;4;1;27\007'`,
		`printf '\033]9;4;1;44\007'`,
		`printf '\033]9;4;1;61\007'`,
		`printf '\033]9;4;1;70\007'`,
	}

	currentIndex := 0
	for _, value := range expectedProgress {
		index := strings.Index(got[currentIndex:], value)
		assert.GreaterOrEqual(t, index, 0, "missing progress command %s", value)
		currentIndex += index + len(value)
	}
}

func TestInitAutoProgressStartsAfterInternalSpan(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `progress:
  start_percentage: 10
  internal: 10
  end_percentage: 70
alias:
  - name: ll
    value: ls -la
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	got := Init(configFile, shell.BASH, true)
	assert.Contains(t, got, `printf '\033]9;4;1;20\007'`)
}

func TestInitAllowsUnknownPropertiesAtRuntime(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `env:
  - name: ANDROID_HOME
    value: '{{ .Home }}/Android/Sdk'
    isPath: true
    mysterySetting: true
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	assert.Contains(t, script, "export ANDROID_HOME=")
}

func TestInitRejectsNonPositiveScriptWeight(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `script:
  - value: echo hello
    weight: 0
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	assert.Contains(t, script, "aliae error:")
	assert.Contains(t, script, "script[0].weight")
	assert.Contains(t, script, "greater than 0")
}
