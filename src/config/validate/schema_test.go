package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeSchemaPath(t *testing.T) {
	assert.Equal(t, "", NormalizeSchemaPath("(root)"))
	assert.Equal(t, "alias.0.name", NormalizeSchemaPath("alias.0.name"))
}

func TestLineResolverAnnotate(t *testing.T) {
	data := []byte("alias:\n  - value: demo\n")

	resolver, err := NewLineResolver(data)
	require.NoError(t, err)

	annotated := resolver.Annotate("alias.0.value", "name is required")
	assert.Contains(t, annotated, "line 2:")
	assert.Contains(t, annotated, "value: demo")
}
