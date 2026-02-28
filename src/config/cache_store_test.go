package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigCacheCreatesAndInvalidates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
alias:
  - name: demo
    value: one
`), 0o600))

	first, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.Len(t, first.Aliae, 1)
	assert.Equal(t, "one", string(first.Aliae[0].Value))
	assert.True(t, first.Cache)

	_, statErr := os.Stat(configCachePath(configFile))
	assert.NoError(t, statErr)

	time.Sleep(2 * time.Millisecond)
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
alias:
  - name: demo
    value: two
`), 0o600))

	second, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.Len(t, second.Aliae, 1)
	assert.Equal(t, "two", string(second.Aliae[0].Value))
	assert.True(t, second.Cache)
}
