package init

import (
	"os"
	"path/filepath"
	"testing"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRenderCacheInvalidatesOnTrackedEnvChange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("RENDER_CACHE_ENV", "one")

	configFile := filepath.ToSlash(filepath.Join(t.TempDir(), "aliae.yaml"))
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
env:
  - name: RENDER_CACHE_ENV
    value: '{{ .Env.RENDER_CACHE_ENV }}'
`), 0o600))

	first := Init(configFile, shell.BASH, true)
	assert.Contains(t, first, `export RENDER_CACHE_ENV="one"`)

	aliae, err := cfg.LoadConfig(configFile)
	require.NoError(t, err)
	cacheState := prepareRenderCache(configFile, shell.BASH, aliae)
	require.NotNil(t, cacheState)
	_, statErr := os.Stat(cacheState.entryPath)
	require.NoError(t, statErr)

	t.Setenv("RENDER_CACHE_ENV", "two")
	second := Init(configFile, shell.BASH, true)
	assert.Contains(t, second, `export RENDER_CACHE_ENV="two"`)
}

func TestInitRenderCacheSkippedForStatefulScripts(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	configFile := filepath.ToSlash(filepath.Join(t.TempDir(), "aliae.yaml"))
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
script:
  - value: echo once
    state:
      file: run-once.state
`), 0o600))

	_ = Init(configFile, shell.BASH, true)

	aliae, err := cfg.LoadConfig(configFile)
	require.NoError(t, err)
	assert.Nil(t, prepareRenderCache(configFile, shell.BASH, aliae))
}
