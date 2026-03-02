package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRenderer struct {
	calls *[]string
	id    string
}

func (t testRenderer) Render() {
	*t.calls = append(*t.calls, t.id)
}

func TestRenderCompositeRenderOrder(t *testing.T) {
	calls := []string{}
	composite := NewRenderComposite(
		testRenderer{id: "env", calls: &calls},
		testRenderer{id: "path", calls: &calls},
		testRenderer{id: "alias", calls: &calls},
	)

	composite.Render()
	assert.Equal(t, []string{"env", "path", "alias"}, calls)
}

func TestRenderCompositeSkipsNilRenderers(t *testing.T) {
	calls := []string{}
	composite := NewRenderComposite(
		testRenderer{id: "env", calls: &calls},
		nil,
		testRenderer{id: "script", calls: &calls},
	)

	composite.Render()
	assert.Equal(t, []string{"env", "script"}, calls)
}
