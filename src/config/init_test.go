package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestInitTemplateConfigVariables(t *testing.T) {
	shell.DotFile.Reset()
	t.Cleanup(shell.DotFile.Reset)

	tempDir := t.TempDir()
	configFile := filepath.ToSlash(filepath.Join(tempDir, "aliae.yaml"))
	configContent := `env:
  - name: ALIAE_CONFIG_PATH
    value: '{{ .AliaeConfig }}'
  - name: ALIAE_CONFIG_DIR
    value: '{{ .AliaeConfigDir }}'
  - name: CONFIG_PATH
    value: '{{ .ConfigPath }}'
  - name: CONFIG_DIR
    value: '{{ .ConfigDir }}'
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	escapedDir := strings.ReplaceAll(resolveConfigDir(configFile), `\`, `\\`)
	expected := "export ALIAE_CONFIG_PATH=\"" + configFile + "\"\n" +
		"export ALIAE_CONFIG_DIR=\"" + escapedDir + "\"\n" +
		"export CONFIG_PATH=\"" + configFile + "\"\n" +
		"export CONFIG_DIR=\"" + escapedDir + "\""
	assert.Equal(t, expected, script)
}
