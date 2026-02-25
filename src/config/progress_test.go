package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigProgressDefaults(t *testing.T) {
	aliae, err := parseConfig([]byte(`alias:
  - name: g
    value: git
`))
	require.NoError(t, err)
	assert.False(t, aliae.Progress.Enabled)
}

func TestParseConfigProgressReset(t *testing.T) {
	aliae, err := parseConfig([]byte(`progress:
  start_percentage: 21.44
  end_percentage: reset
`))
	require.NoError(t, err)

	assert.True(t, aliae.Progress.Enabled)
	assert.Equal(t, 21.44, aliae.Progress.StartPercentage)
	assert.True(t, aliae.Progress.EndPercentage.Reset)
	assert.Equal(t, 100.0, aliae.Progress.EndPercentage.Value)
}

func TestParseConfigProgressFalse(t *testing.T) {
	aliae, err := parseConfig([]byte(`progress: false`))
	require.NoError(t, err)
	assert.False(t, aliae.Progress.Enabled)
}

func TestParseConfigProgressTrueUsesDefaults(t *testing.T) {
	aliae, err := parseConfig([]byte(`progress: true`))
	require.NoError(t, err)
	assert.True(t, aliae.Progress.Enabled)
	assert.Equal(t, 0.0, aliae.Progress.StartPercentage)
	assert.True(t, aliae.Progress.EndPercentage.Reset)
	assert.Equal(t, 100.0, aliae.Progress.EndPercentage.Value)
}

func TestProgressUnmarshalYAMLMethod(t *testing.T) {
	var doc struct {
		Progress Progress `yaml:"progress"`
	}

	err := yaml.Unmarshal([]byte(`progress:
  start_percentage: 10
  end_percentage: 70
`), &doc)
	require.NoError(t, err)
	assert.True(t, doc.Progress.Enabled)
	assert.Equal(t, 10.0, doc.Progress.StartPercentage)
	assert.Equal(t, 70.0, doc.Progress.EndPercentage.Value)
}

func TestLoadConfigProgressObject(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(file, []byte(`progress:
  start_percentage: 10
  end_percentage: 70
`), 0o600))

	aliae, err := LoadConfig(file)
	require.NoError(t, err)
	assert.True(t, aliae.Progress.Enabled)
	assert.Equal(t, 10.0, aliae.Progress.StartPercentage)
	assert.Equal(t, 70.0, aliae.Progress.EndPercentage.Value)
	assert.False(t, aliae.Progress.EndPercentage.Reset)
}
