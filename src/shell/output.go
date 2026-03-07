package shell

import (
	"strings"
	"sync"
)

var (
	renderOutputMu sync.RWMutex
	defaultOutput  strings.Builder
	renderOutput   = &defaultOutput
)

func SetRenderOutput(builder *strings.Builder) func() {
	if builder == nil {
		builder = &defaultOutput
	}

	renderOutputMu.Lock()
	previous := renderOutput
	renderOutput = builder
	renderOutputMu.Unlock()

	return func() {
		renderOutputMu.Lock()
		renderOutput = previous
		renderOutputMu.Unlock()
	}
}

func ResetRenderOutput() {
	activeRenderOutput().Reset()
}

func RenderOutputString() string {
	return activeRenderOutput().String()
}

func WriteRenderOutput(value string) {
	writeRenderOutput(value)
}

func renderOutputLen() int {
	return activeRenderOutput().Len()
}

func writeRenderOutput(value string) {
	activeRenderOutput().WriteString(value)
}

func activeRenderOutput() *strings.Builder {
	renderOutputMu.RLock()
	defer renderOutputMu.RUnlock()
	return renderOutput
}
