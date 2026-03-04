package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubRenderer struct {
	id    string
	calls *[]string
}

func (s stubRenderer) Render() {
	*s.calls = append(*s.calls, s.id)
}

func TestRenderCompositeRenderOrder(t *testing.T) {
	calls := []string{}
	composite := newRenderComposite(
		stubRenderer{id: "env", calls: &calls},
		stubRenderer{id: "path", calls: &calls},
		stubRenderer{id: "alias", calls: &calls},
	)

	composite.Render()

	assert.Equal(t, []string{"env", "path", "alias"}, calls)
}

func TestRenderCompositeSkipsNilRenderers(t *testing.T) {
	calls := []string{}
	composite := newRenderComposite(
		nil,
		stubRenderer{id: "env", calls: &calls},
		nil,
	)

	composite.Render()

	assert.Equal(t, []string{"env"}, calls)
}
