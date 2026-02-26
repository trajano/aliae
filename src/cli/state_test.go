package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunStateListAndClear(t *testing.T) {
	originalConfig := config
	originalRuntime := context.Current
	t.Cleanup(func() {
		config = originalConfig
		context.Current = originalRuntime
	})

	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`script:
  - value: echo hello
    state:
      file: hello.state
      runEvery: 24h
  - value: echo bye
    state:
      file: bye.state
`), 0o600))

	context.Current = &context.Runtime{
		OS:   context.LINUX,
		Home: root,
	}
	config = configFile

	helloStatePath := filepath.Join(root, ".local", "aliae", "state", "hello.state")
	require.NoError(t, aliaeState.WriteLastRun(helloStatePath, time.Date(2026, 2, 25, 10, 0, 0, 0, time.UTC), aliaeState.FormatJSON))

	var listOut bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&listOut)
	cmd.SetErr(&bytes.Buffer{})
	require.NoError(t, runStateList(cmd))
	assert.Contains(t, listOut.String(), "FILE")
	assert.Contains(t, listOut.String(), "hello.state")
	assert.Contains(t, listOut.String(), "bye.state")
	assert.Contains(t, listOut.String(), "24h0m0s")
	assert.Contains(t, listOut.String(), "once")

	byeStatePath := filepath.Join(root, ".local", "aliae", "state", "bye.state")
	require.NoError(t, aliaeState.WriteLastRun(byeStatePath, time.Now(), aliaeState.FormatText))

	var clearOut bytes.Buffer
	clearCmd := &cobra.Command{}
	clearCmd.SetOut(&clearOut)
	clearCmd.SetErr(&bytes.Buffer{})
	require.NoError(t, runStateClear(clearCmd))

	_, err := os.Stat(helloStatePath)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = os.Stat(byeStatePath)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Contains(t, clearOut.String(), "removed hello.state")
	assert.Contains(t, clearOut.String(), "removed bye.state")
}

func TestRunStateListNoEntries(t *testing.T) {
	originalConfig := config
	originalRuntime := context.Current
	t.Cleanup(func() {
		config = originalConfig
		context.Current = originalRuntime
	})

	root := t.TempDir()
	configFile := filepath.Join(root, "aliae.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(`script:
  - value: echo hello
`), 0o600))

	context.Current = &context.Runtime{
		OS:   context.LINUX,
		Home: root,
	}
	config = configFile

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	require.NoError(t, runStateList(cmd))
	assert.Contains(t, out.String(), "No state entries referenced in config.")
}
