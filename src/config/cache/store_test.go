package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAndLoadRoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	source := filepath.Join(root, "config.yaml")
	require.NoError(t, os.WriteFile(source, []byte("alias: []\n"), 0o600))

	cachePath := Path(source, "shell=bash;os=linux;wsl=false")
	payload := []byte("cached-config")

	require.NoError(t, Store(cachePath, "schema-hash", true, []string{source}, payload))

	loaded, ok, err := Load(cachePath, "schema-hash", true)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, payload, loaded)
}

func TestLoadInvalidatesOnSourceChange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	source := filepath.Join(root, "config.yaml")
	require.NoError(t, os.WriteFile(source, []byte("alias:\n  - name: one\n"), 0o600))

	cachePath := Path(source, "shell=bash;os=linux;wsl=false")
	require.NoError(t, Store(cachePath, "schema-hash", false, []string{source}, []byte("payload")))

	time.Sleep(2 * time.Millisecond)
	require.NoError(t, os.WriteFile(source, []byte("alias:\n  - name: two\n"), 0o600))

	loaded, ok, err := Load(cachePath, "schema-hash", false)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, loaded)
}

func TestPathIsDeterministicAndContextSensitive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	configPath := filepath.Join(t.TempDir(), "aliae.yaml")
	pathA := Path(configPath, "shell=bash;os=linux;wsl=false")
	pathB := Path(configPath, "shell=bash;os=linux;wsl=false")
	pathC := Path(configPath, "shell=pwsh;os=windows;wsl=false")

	assert.Equal(t, pathA, pathB)
	assert.NotEqual(t, pathA, pathC)
}

func TestStateFlags(t *testing.T) {
	SetBypass(false)
	assert.False(t, IsBypassed())
	SetBypass(true)
	assert.True(t, IsBypassed())

	SetLastLoadUsed(false)
	assert.False(t, LastLoadUsed())
	SetLastLoadUsed(true)
	assert.True(t, LastLoadUsed())
}
