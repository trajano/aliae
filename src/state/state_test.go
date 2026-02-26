package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootDirLinux(t *testing.T) {
	context.Current = &context.Runtime{
		OS:   context.LINUX,
		Home: "/home/tester",
	}

	assert.Equal(t, filepath.Join("/home/tester", ".local", "aliae", "state"), RootDir())
}

func TestRootDirDarwin(t *testing.T) {
	context.Current = &context.Runtime{
		OS:   context.DARWIN,
		Home: "/Users/tester",
	}

	assert.Equal(t, filepath.Join("/Users/tester", "Library", "Application Support", "aliae", "State"), RootDir())
}

func TestRootDirWindows(t *testing.T) {
	context.Current = &context.Runtime{
		OS:   context.WINDOWS,
		Home: `C:\Users\tester`,
	}
	t.Setenv("LOCALAPPDATA", `C:\Users\tester\AppData\Local`)

	assert.Equal(t, filepath.Join(`C:\Users\tester\AppData\Local`, "aliae", "state"), RootDir())
}

func TestShouldRunAndWriteLastRun(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, "script.json")
	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)

	run, last, err := ShouldRun(statePath, 0, now)
	require.NoError(t, err)
	assert.True(t, run)
	assert.Nil(t, last)

	require.NoError(t, WriteLastRun(statePath, now, FormatJSON))

	run, last, err = ShouldRun(statePath, 0, now.Add(time.Minute))
	require.NoError(t, err)
	assert.False(t, run)
	require.NotNil(t, last)
	assert.Equal(t, now, *last)

	run, _, err = ShouldRun(statePath, 2*time.Hour, now.Add(time.Hour))
	require.NoError(t, err)
	assert.False(t, run)

	run, _, err = ShouldRun(statePath, 2*time.Hour, now.Add(3*time.Hour))
	require.NoError(t, err)
	assert.True(t, run)
}

func TestReadLastRunText(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, "script.txt")
	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)

	require.NoError(t, os.WriteFile(statePath, []byte(now.Format(time.RFC3339Nano)), 0o600))

	got, err := ReadLastRun(statePath)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, now, *got)
}

func TestIsValidFileName(t *testing.T) {
	assert.True(t, IsValidFileName("my-script.state"))
	assert.False(t, IsValidFileName(""))
	assert.False(t, IsValidFileName("../state"))
	assert.False(t, IsValidFileName("nested/state"))
}
