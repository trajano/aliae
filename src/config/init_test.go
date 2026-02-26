package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestInitInternalProgressSequence(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	base1 := filepath.ToSlash(filepath.Join(tempDir, "base1.yaml"))
	base2 := filepath.ToSlash(filepath.Join(tempDir, "base2.yaml"))
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))

	require.NoError(t, os.WriteFile(base1, []byte(`alias:
  - name: b1
    value: one
`), 0o600))
	require.NoError(t, os.WriteFile(base2, []byte(`alias:
  - name: b2
    value: two
`), 0o600))
	require.NoError(t, os.WriteFile(configFile, []byte(`progress:
  start_percentage: 10
  internal: 10
  end_percentage: reset
extends:
  - ./base1.yaml
  - ./base2.yaml
alias:
  - name: root
    value: root
`), 0o600))

	var stderr bytes.Buffer
	SetInitProgressWriter(&stderr)
	t.Cleanup(resetInitProgressWriter)

	_ = Init(configFile, shell.BASH, true)

	output := stderr.String()
	expected := []string{
		"\x1b]9;4;1;10\x07",
		"\x1b]9;4;1;11\x07",
		"\x1b]9;4;1;12\x07",
		"\x1b]9;4;1;13\x07",
		"\x1b]9;4;1;14\x07",
		"\x1b]9;4;1;16\x07",
		"\x1b]9;4;1;17\x07",
	}

	currentIndex := 0
	for _, value := range expected {
		index := strings.Index(output[currentIndex:], value)
		assert.GreaterOrEqual(t, index, 0, "missing internal progress %q", value)
		currentIndex += index + len(value)
	}
}

func TestInitInternalProgressIncludesStateChecks(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))

	require.NoError(t, os.WriteFile(configFile, []byte(`progress:
  start_percentage: 10
  internal: 5
  end_percentage: reset
script:
  - value: echo once
    state:
      file: once.state
`), 0o600))

	var stderr bytes.Buffer
	SetInitProgressWriter(&stderr)
	t.Cleanup(resetInitProgressWriter)

	_ = Init(configFile, shell.BASH, true)

	output := stderr.String()
	assert.Contains(t, output, "\x1b]9;4;1;13\x07")
	assert.NotContains(t, output, "\x1b]9;4;1;14\x07")
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

func TestInitRejectsInvalidScriptStateRunEvery(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `script:
  - value: echo hello
    state:
      file: hello.state
      runEvery: nope
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	assert.Contains(t, script, "aliae error:")
	assert.Contains(t, script, "script[0].state.runEvery")
}
