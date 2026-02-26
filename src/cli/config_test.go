package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderResolvedConfigYAML(t *testing.T) {
	root := t.TempDir()
	aliases := filepath.Join(root, "aliases")
	fragments := filepath.Join(root, "fragments")
	require.NoError(t, os.MkdirAll(aliases, 0o700))
	require.NoError(t, os.MkdirAll(fragments, 0o700))

	require.NoError(t, os.WriteFile(filepath.Join(fragments, "base.yaml"), []byte(`alias:
  - name: base
    value: from-base
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(aliases, "10-nested.yaml"), []byte(`extends:
  - ../fragments/base.yaml
alias:
  - name: nested
    value: from-nested
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "env.yaml"), []byte(`- name: TEST_ENV
  value: test
`), 0o600))

	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`extends:
  - dir: ./aliases
env: !include ./env.yaml
script:
  - value: echo hello
    state:
      file: hello.state
      runEvery: 24h
alias:
  - name: root
    value: from-root
    type: command
`), 0o600))

	output, err := renderResolvedConfigYAML(configFile)
	require.NoError(t, err)

	assert.Contains(t, output, "alias:")
	assert.Contains(t, output, "name: base")
	assert.Contains(t, output, "name: nested")
	assert.Contains(t, output, "name: root")
	assert.Contains(t, output, "env:")
	assert.Contains(t, output, "name: TEST_ENV")
	assert.Contains(t, output, "script:")
	assert.Contains(t, output, "state:")
	assert.Contains(t, output, "file: hello.state")
	assert.Contains(t, output, "runEvery: 24h")
	assert.NotContains(t, output, "weight: 0")
	assert.NotContains(t, output, "description: \"\"")
	assert.NotContains(t, output, "type: \"\"")
	assert.NotContains(t, output, "force: false")
	assert.NotContains(t, output, "extends:")
	assert.NotContains(t, output, "!include")

	rootName := "name: root"
	rootType := "type: command"
	rootValue := "value: from-root"
	nameIndex := strings.Index(output, rootName)
	typeIndex := strings.Index(output, rootType)
	valueIndex := strings.Index(output, rootValue)
	assert.Greater(t, nameIndex, -1)
	assert.Greater(t, typeIndex, -1)
	assert.Greater(t, valueIndex, -1)
	assert.Less(t, nameIndex, typeIndex)
	assert.Less(t, typeIndex, valueIndex)
}
