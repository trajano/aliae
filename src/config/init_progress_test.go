package config

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestProgressTotalWeightIncludesVarCDPathAndLink(t *testing.T) {
	cfg := &Aliae{
		Vars: Vars{
			{Name: "A"},
		},
		Envs: shell.Envs{
			{Name: "E"},
		},
		Paths: shell.Paths{
			{Value: "/bin"},
		},
		CDPaths: shell.CDPaths{
			{Value: "/tmp"},
		},
		Aliae: shell.Aliae{
			{Name: "a"},
		},
		Links: shell.Links{
			{Name: "l", Target: "t"},
		},
		Scripts: shell.Scripts{
			{Value: "echo one"},
			{Value: "echo two", Weight: 2},
		},
	}

	assert.Equal(t, 9.0, cfg.progressTotalWeight())
}
