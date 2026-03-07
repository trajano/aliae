package shell

import (
	"sync"

	"github.com/jandedobbeleer/aliae/src/context"
)

var (
	activeRuntimeMu sync.RWMutex
	activeRuntime   *context.Runtime
)

func SetRuntime(runtime *context.Runtime) func() {
	activeRuntimeMu.Lock()
	previous := activeRuntime
	activeRuntime = runtime
	activeRuntimeMu.Unlock()

	return func() {
		activeRuntimeMu.Lock()
		activeRuntime = previous
		activeRuntimeMu.Unlock()
	}
}

func currentRuntime() *context.Runtime {
	activeRuntimeMu.RLock()
	current := activeRuntime
	activeRuntimeMu.RUnlock()
	if current != nil {
		return current
	}

	return context.Current
}
