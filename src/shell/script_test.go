package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptRender(t *testing.T) {
	cases := []struct {
		Case           string
		Expected       string
		Scripts        Scripts
		NonEmptyScript bool
	}{
		{
			Case:    "No content",
			Scripts: Scripts{},
		},
		{
			Case: "Simple script",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			Expected: "foo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Non-Empty",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			NonEmptyScript: true,
			Expected:       "foo\n\nfoo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Multiple scripts",
			Scripts: Scripts{
				{
					Value: "foo",
				},
				{
					Value: "bar",
				},
			},
			Expected: "foo\nbar",
		},
		{
			Case: "Python script type",
			Scripts: Scripts{
				{
					Type:  PythonScript,
					Value: "print(123)",
				},
			},
			Expected: `python -c "print(123)"`,
		},
		{
			Case: "Perl script type",
			Scripts: Scripts{
				{
					Type:  PerlScript,
					Value: "print qq(123)",
				},
			},
			Expected: `perl -e "print qq(123)"`,
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: PWSH}
		tc.Scripts.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}

func TestScriptRenderCmdPythonAndPerl(t *testing.T) {
	cases := []struct {
		Case     string
		Script   *Script
		Expected string
	}{
		{
			Case: "Python",
			Script: &Script{
				Type:  PythonScript,
				Value: "print(123)",
			},
			Expected: `os.execute("python -c \"print(123)\"")`,
		},
		{
			Case: "Perl",
			Script: &Script{
				Type:  PerlScript,
				Value: "print qq(123)",
			},
			Expected: `os.execute("perl -e \"print qq(123)\"")`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: CMD}
		assert.Equal(t, tc.Expected, tc.Script.String(), tc.Case)
	}
}

func TestScriptStateRunOnce(t *testing.T) {
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: PWSH, OS: context.LINUX, Home: tempDir}

	scripts := Scripts{
		{
			Value: "echo hello",
			State: ScriptState{
				File: "run-once.state",
			},
		},
	}

	DotFile.Reset()
	scripts.Render()
	assert.Equal(t, "echo hello", strings.TrimSpace(DotFile.String()))

	DotFile.Reset()
	scripts.Render()
	assert.Equal(t, "", strings.TrimSpace(DotFile.String()))
}

func TestScriptStateRunEvery(t *testing.T) {
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: PWSH, OS: context.LINUX, Home: tempDir}

	statePath := filepath.Join(tempDir, ".local", "aliae", "state", "hourly.state")
	old := time.Now().Add(-2 * time.Hour)
	require.NoError(t, aliaeState.WriteLastRun(statePath, old, aliaeState.FormatJSON))

	scripts := Scripts{
		{
			Value: "echo hello",
			State: ScriptState{
				File:     "hourly.state",
				RunEvery: "1h",
			},
		},
	}

	DotFile.Reset()
	scripts.Render()
	assert.Equal(t, "echo hello", strings.TrimSpace(DotFile.String()))

	updatedLastRun, err := aliaeState.ReadLastRun(statePath)
	require.NoError(t, err)
	require.NotNil(t, updatedLastRun)
	assert.True(t, updatedLastRun.After(old))
}

func TestScriptStateRunEveryNotDue(t *testing.T) {
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: PWSH, OS: context.LINUX, Home: tempDir}

	statePath := filepath.Join(tempDir, ".local", "aliae", "state", "daily.state")
	require.NoError(t, os.MkdirAll(filepath.Dir(statePath), 0o700))
	require.NoError(t, os.WriteFile(statePath, []byte(time.Now().Format(time.RFC3339Nano)), 0o600))

	scripts := Scripts{
		{
			Value: "echo hello",
			State: ScriptState{
				File:     "daily.state",
				RunEvery: "24h",
			},
		},
	}

	DotFile.Reset()
	scripts.Render()
	assert.Equal(t, "", strings.TrimSpace(DotFile.String()))
}
