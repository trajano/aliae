package shell

import (
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestAutoProgressRendersAndResets(t *testing.T) {
	DotFile.Reset()
	context.Current = &context.Runtime{Shell: BASH}

	StartAutoProgress(AutoProgressConfig{
		Enabled:         true,
		StartPercentage: 20,
		EndPercentage:   100,
		ResetAtEnd:      true,
		TotalWeight:     4,
	})

	advanceAutoProgress(1)
	advanceAutoProgress(1)
	advanceAutoProgress(1)
	advanceAutoProgress(1)
	EndAutoProgress()

	got := strings.Split(strings.TrimSpace(DotFile.String()), "\n")
	assert.Equal(
		t,
		[]string{
			`printf '\033]9;4;1;20\007'`,
			`printf '\033]9;4;1;40\007'`,
			`printf '\033]9;4;1;60\007'`,
			`printf '\033]9;4;1;80\007'`,
			`printf '\033]9;4;1;99\007'`,
			`printf '\033]9;4;0;0\007'`,
		},
		got,
	)
}

func TestAutoProgressDisabled(t *testing.T) {
	DotFile.Reset()
	context.Current = &context.Runtime{Shell: BASH}

	StartAutoProgress(AutoProgressConfig{
		Enabled:     false,
		TotalWeight: 10,
	})
	advanceAutoProgress(1)
	EndAutoProgress()

	assert.Equal(t, "", DotFile.String())
}
