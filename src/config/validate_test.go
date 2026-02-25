package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`alias:
  - name: g
    value: git
`), 0o600))

		err := ValidateConfig(file)
		require.NoError(t, err)
	})

	t.Run("invalid config", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`alias: invalid
`), 0o600))

		err := ValidateConfig(file)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "failed to parse config file")
	})

	t.Run("schema error includes source line", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`alias:
  - name: g
    value: git
    type: nope
`), 0o600))

		err := ValidateConfig(file)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "schema validation failed")
		assert.Contains(t, strings.ToLower(err.Error()), "line")
		assert.Contains(t, err.Error(), "type: nope")
	})

	t.Run("include tags are validated after resolution", func(t *testing.T) {
		root := t.TempDir()
		includeFile := filepath.Join(root, "aliases.yaml")
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(includeFile, []byte(`- name: g
  value: git
`), 0o600))
		require.NoError(t, os.WriteFile(file, []byte(`alias: !include ./aliases.yaml
`), 0o600))

		err := ValidateConfig(file)
		require.NoError(t, err)
	})

	t.Run("invalid if expression", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`alias:
  - name: g
    value: git
    if: '{'
`), 0o600))

		err := ValidateConfig(file)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "if expression validation failed")
		assert.Contains(t, err.Error(), "alias[0].if")
		assert.Contains(t, strings.ToLower(err.Error()), "line")
		assert.Contains(t, err.Error(), `if: "{"`)
	})

	t.Run("invalid progress percentages", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`progress:
  start_percentage: 60
  end_percentage: 20
`), 0o600))

		err := ValidateConfig(file)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "progress validation failed")
		assert.Contains(t, err.Error(), "progress.end_percentage")
	})

	t.Run("invalid script weight", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "aliae.yaml")
		require.NoError(t, os.WriteFile(file, []byte(`script:
  - value: echo hello
    weight: 0.5
`), 0o600))

		err := ValidateConfig(file)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "schema validation failed")
		assert.Contains(t, err.Error(), "weight")
	})
}
