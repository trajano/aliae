package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetShellWritesToStdout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := getCmd.RunE(cmd, []string{"shell"})
	require.NoError(t, err)
	assert.NotEmpty(t, strings.TrimSpace(stdout.String()))
	assert.Empty(t, stderr.String())
}

func TestGetBenchmarkWithoutShell(t *testing.T) {
	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`alias:
  - name: g
    value: git
`), 0o600))

	previousConfig := config
	config = configFile
	t.Cleanup(func() {
		config = previousConfig
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := getCmd.RunE(cmd, []string{"benchmark"})
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "aliae get benchmark")
	assert.Contains(t, stdout.String(), "benchmark.cache_used=")
	assert.Contains(t, stdout.String(), "benchmark.cygpath=")
	assert.Contains(t, stdout.String(), "benchmark.load_config=")
	assert.Contains(t, stdout.String(), "benchmark.evaluate_vars=")
	assert.Contains(t, stdout.String(), "benchmark.render_config=")
	assert.Contains(t, stdout.String(), "benchmark.validate_config=skipped")
	assert.Contains(t, stdout.String(), "benchmark.total=")
	assert.Empty(t, stderr.String())
}

func TestGetBenchmarkWithShell(t *testing.T) {
	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`alias:
  - name: g
    value: git
`), 0o600))

	previousConfig := config
	config = configFile
	t.Cleanup(func() {
		config = previousConfig
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := getCmd.RunE(cmd, []string{"benchmark", "bash"})
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "benchmark.shell=bash")
	assert.Contains(t, stdout.String(), "benchmark.cache_used=")
	assert.Contains(t, stdout.String(), "benchmark.cygpath=")
	assert.Contains(t, stdout.String(), "benchmark.generate_init_bash=")
	assert.Contains(t, stdout.String(), "benchmark.init.render_env=")
	assert.Contains(t, stdout.String(), "benchmark.init.render_script=")
	assert.Empty(t, stderr.String())
}

func TestGetBenchmarkNoCacheFlagDisablesCacheUsage(t *testing.T) {
	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`cache: true
alias:
  - name: g
    value: git
`), 0o600))

	previousConfig := config
	previousNoCache := benchmarkNoCache
	config = configFile
	benchmarkNoCache = true
	t.Cleanup(func() {
		config = previousConfig
		benchmarkNoCache = previousNoCache
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := getCmd.RunE(cmd, []string{"benchmark"})
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "benchmark.cache_used=false")
	assert.Contains(t, stdout.String(), "benchmark.cygpath=")
	assert.Contains(t, stdout.String(), "benchmark.validate_config=")
	assert.NotContains(t, stdout.String(), "benchmark.validate_config=skipped")
	assert.Empty(t, stderr.String())
}

func TestValidateGetArgsBenchmarkShellValidation(t *testing.T) {
	err := validateGetArgs(getCmd, []string{"benchmark", "not-a-shell"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}

func TestGetBenchmarkReturnsErrorWhenConfigLoadFails(t *testing.T) {
	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`extends:
  - path: does-not-exist.yml
`), 0o600))

	previousConfig := config
	config = configFile
	t.Cleanup(func() {
		config = previousConfig
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	require.NotPanics(t, func() {
		err := getCmd.RunE(cmd, []string{"benchmark", "bash"})
		require.Error(t, err)
	})
}
