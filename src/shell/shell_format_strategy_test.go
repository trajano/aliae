package shell

import (
	"reflect"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestFormatStrategySelection(t *testing.T) {
	cases := []struct {
		shell        string
		expectedType string
	}{
		{shell: ZSH, expectedType: "shell.zshFormatStrategy"},
		{shell: BASH, expectedType: "shell.bashFormatStrategy"},
		{shell: PWSH, expectedType: "shell.pwshFormatStrategy"},
		{shell: POWERSHELL, expectedType: "shell.pwshFormatStrategy"},
		{shell: NU, expectedType: "shell.nuFormatStrategy"},
		{shell: FISH, expectedType: "shell.fishFormatStrategy"},
		{shell: TCSH, expectedType: "shell.tcshFormatStrategy"},
		{shell: XONSH, expectedType: "shell.xonshFormatStrategy"},
		{shell: CMD, expectedType: "shell.cmdFormatStrategy"},
		{shell: "unknown", expectedType: "shell.noopFormatStrategy"},
	}

	for _, tc := range cases {
		useRuntime(t, &context.Runtime{Shell: tc.shell})
		assert.Equal(t, tc.expectedType, reflect.TypeOf(formatStrategy()).String(), tc.shell)
	}

	useRuntime(t, nil)
	assert.Equal(t, "shell.noopFormatStrategy", reflect.TypeOf(formatStrategy()).String())
}

func TestCDPathCurrentDirScriptByStrategy(t *testing.T) {
	cases := []struct {
		shell    string
		expected string
	}{
		{
			shell:    BASH,
			expected: `if [ -n "$CDPATH" ]; then export CDPATH=":$CDPATH"; else export CDPATH=":"; fi`,
		},
		{shell: ZSH, expected: `cdpath=( . $cdpath )`},
		{shell: FISH, expected: `set -g cdpath . $cdpath`},
		{shell: TCSH, expected: `set cdpath = ( . $cdpath );`},
		{shell: XONSH, expected: `$CDPATH = ["."] + $CDPATH`},
		{shell: PWSH, expected: ""},
		{shell: NU, expected: ""},
		{shell: CMD, expected: ""},
	}

	for _, tc := range cases {
		useRuntime(t, &context.Runtime{Shell: tc.shell})
		assert.Equal(t, tc.expected, formatStrategy().FormatCDPathCurrentDirScript(), tc.shell)
	}
}
