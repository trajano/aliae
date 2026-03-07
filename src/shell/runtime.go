package shell

import (
	"sync/atomic"

	"github.com/jandedobbeleer/aliae/src/context"
)

var (
	activeRuntime atomic.Pointer[context.Runtime]
)

func SetRuntime(runtime *context.Runtime) func() {
	previous := activeRuntime.Swap(runtime)

	return func() {
		activeRuntime.Store(previous)
	}
}

func currentRuntime() *context.Runtime {
	return activeRuntime.Load()
}
