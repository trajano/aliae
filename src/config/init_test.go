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
  - name: CONFIG_PATH
    value: '{{ .ConfigPath }}'
  - name: CONFIG_DIR
    value: '{{ .ConfigDir }}'
`

	err := os.WriteFile(configFile, []byte(configContent), 0o600)
	assert.NoError(t, err)

	script := Init(configFile, shell.BASH, true)
	escapedDir := strings.ReplaceAll(resolveConfigDir(configFile), `\`, `\\`)
	expected := "export CONFIG_PATH=\"" + configFile + "\"\n" +
		"export CONFIG_DIR=\"" + escapedDir + "\""
	assert.Equal(t, expected, script)
}
