package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestTemplateParseWithRuntimeUsesProvidedShell(t *testing.T) {
	useRuntime(t, &context.Runtime{Shell: BASH})

	template := Template(`{{ setArg "value" 1 }}`)
	runtime := &context.Runtime{Shell: PWSH}

	parsed := template.ParseWithRuntime(runtime)
	assert.Equal(t, `$value = $args[0]`, string(parsed))
}
