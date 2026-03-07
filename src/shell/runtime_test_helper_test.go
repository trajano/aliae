package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
)

func useRuntime(t *testing.T, runtime *context.Runtime) {
	t.Helper()
	context.Current = runtime
	restore := SetRuntime(runtime)
	t.Cleanup(restore)
}
