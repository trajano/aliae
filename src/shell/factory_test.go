package shell

import (
	"reflect"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestRenderStrategySelection(t *testing.T) {
	cases := []struct {
		shell        string
		expectedType string
	}{
		{shell: ZSH, expectedType: "shell.zshRenderStrategy"},
		{shell: BASH, expectedType: "shell.bashRenderStrategy"},
		{shell: PWSH, expectedType: "shell.pwshRenderStrategy"},
		{shell: POWERSHELL, expectedType: "shell.pwshRenderStrategy"},
		{shell: NU, expectedType: "shell.nuRenderStrategy"},
		{shell: FISH, expectedType: "shell.fishRenderStrategy"},
		{shell: TCSH, expectedType: "shell.tcshRenderStrategy"},
		{shell: XONSH, expectedType: "shell.xonshRenderStrategy"},
		{shell: CMD, expectedType: "shell.cmdRenderStrategy"},
		{shell: "unknown", expectedType: "shell.noopRenderStrategy"},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.shell}
		assert.Equal(t, tc.expectedType, reflect.TypeOf(renderStrategy()).String(), tc.shell)
	}

	context.Current = nil
	assert.Equal(t, "shell.noopRenderStrategy", reflect.TypeOf(renderStrategy()).String())
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
		context.Current = &context.Runtime{Shell: tc.shell}
		assert.Equal(t, tc.expected, renderStrategy().renderCDPathCurrentDirScript(), tc.shell)
	}
}
