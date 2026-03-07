package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
)

func useRuntime(t *testing.T, runtime *context.Runtime) {
	t.Helper()
	context.SetCurrent(runtime)
	restore := SetRuntime(runtime)
	t.Cleanup(func() {
		context.SetCurrent(nil)
		restore()
	})
}
