package shell

import (
	"strings"
)

var (
	defaultOutput strings.Builder
	renderOutput  = &defaultOutput
)

func SetRenderOutput(builder *strings.Builder) func() {
	if builder == nil {
		builder = &defaultOutput
	}

	previous := renderOutput
	renderOutput = builder

	return func() {
		renderOutput = previous
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
	return renderOutput
}
