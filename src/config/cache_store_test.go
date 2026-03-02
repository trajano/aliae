package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/shell"
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

func TestLoadConfigCacheReevaluatesConditionalExtends(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	base := filepath.Join(root, "base.yaml")
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(base, []byte(`alias:
  - name: base
    value: base
`), 0o600))
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
extends:
  - path: ./base.yaml
    if: hasCommandNoCache "aliae-cache-volatile-cmd"
alias:
  - name: child
    value: child
`), 0o600))

	first := Init(configFile, shell.BASH, true)
	assert.NotContains(t, first, `alias base="base"`)
	assert.Contains(t, first, `alias child="child"`)

	originalPath := os.Getenv("PATH")
	tempDir := t.TempDir()
	commandFile := filepath.Join(tempDir, "aliae-cache-volatile-cmd")
	content := []byte("#!/bin/sh\nexit 0\n")
	if runtime.GOOS == "windows" {
		commandFile += ".cmd"
		content = []byte("@echo off\r\nexit /b 0\r\n")
	}
	require.NoError(t, os.WriteFile(commandFile, content, 0o700))
	require.NoError(t, os.Setenv("PATH", tempDir+string(os.PathListSeparator)+originalPath))

	second := Init(configFile, shell.BASH, true)
	assert.Contains(t, second, `alias base="base"`)
	assert.Contains(t, second, `alias child="child"`)
}

func TestLoadConfigCacheRecomputesVarsOnCacheHit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
var:
  - name: ENABLED
    value: true
env:
  - name: CACHE_VAR_TEST
    value: ok
    if: .Var.ENABLED
`), 0o600))

	first := Init(configFile, shell.BASH, true)
	assert.Contains(t, first, `export CACHE_VAR_TEST="ok"`)
	assert.False(t, LastLoadUsedCache())

	shell.DotFile.Reset()
	second := Init(configFile, shell.BASH, true)
	assert.Contains(t, second, `export CACHE_VAR_TEST="ok"`)
	assert.True(t, LastLoadUsedCache())
}

func TestLoadConfigCacheSeparatesConditionalExtendsByShell(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	root := t.TempDir()
	bashConfig := filepath.Join(root, "bash.yml")
	pwshConfig := filepath.Join(root, "pwsh.yml")
	configFile := filepath.Join(root, "aliae.yaml")

	require.NoError(t, os.WriteFile(bashConfig, []byte(`script:
  - value: echo bash-only
`), 0o600))
	require.NoError(t, os.WriteFile(pwshConfig, []byte(`script:
  - value: Write-Host pwsh-only
`), 0o600))
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
extends:
  - path: ./bash.yml
    if: eq .Shell "bash"
  - path: ./pwsh.yml
    if: eq .Shell "pwsh"
`), 0o600))

	bashInit := Init(configFile, shell.BASH, true)
	assert.Contains(t, bashInit, "echo bash-only")
	assert.NotContains(t, bashInit, "Write-Host pwsh-only")
	assert.False(t, LastLoadUsedCache())

	pwshInit := Init(configFile, shell.PWSH, true)
	assert.Contains(t, pwshInit, "Write-Host pwsh-only")
	assert.NotContains(t, pwshInit, "echo bash-only")
	assert.False(t, LastLoadUsedCache())
}
