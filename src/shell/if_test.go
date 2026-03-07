package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	cases := []struct {
		Case     string
		If       If
		Expected bool
	}{
		{
			Case:     "Empty if",
			Expected: false,
		},
		{
			Case:     "Broken if",
			If:       "{}",
			Expected: false,
		},
		{
			Case:     "Match shell",
			If:       `eq .Shell "zsh"`,
			Expected: false,
		},
		{
			Case:     "Hide in current shell",
			If:       `eq .Shell "pwsh"`,
			Expected: true,
		},
		{
			Case:     "Only two shells",
			If:       `match .Shell "bash" "zsh"`,
			Expected: false,
		},
		{
			Case:     "Only two shells",
			If:       `match .Shell "pwsh" "nu"`,
			Expected: true,
		},
	}

	for _, tc := range cases {
		useRuntime(t, &context.Runtime{Shell: "zsh"})
		assert.Equal(t, tc.Expected, tc.If.Ignore(), tc.Case)
	}
}

func TestIfValidate(t *testing.T) {
	t.Run("valid expression", func(t *testing.T) {
		useRuntime(t, &context.Runtime{Shell: "zsh", OS: context.LINUX, Home: "/tmp"})
		err := If(`eq .Shell "zsh"`).Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid expression", func(t *testing.T) {
		useRuntime(t, &context.Runtime{Shell: "zsh", OS: context.LINUX, Home: "/tmp"})
		err := If("{").Validate()
		assert.Error(t, err)
	})

	t.Run("nil runtime context", func(t *testing.T) {
		useRuntime(t, nil)
		err := If(`eq .OS "linux"`).Validate()
		assert.NoError(t, err)
	})
}
