package config

// renderer is the shared contract for init section renderers.
type renderer interface {
	Render()
}

// renderComposite composes multiple renderers and executes them in order.
type renderComposite struct {
	renderers []renderer
}

func newRenderComposite(renderers ...renderer) renderComposite {
	filtered := make([]renderer, 0, len(renderers))
	for _, item := range renderers {
		if item == nil {
			continue
		}

		filtered = append(filtered, item)
	}

	return renderComposite{renderers: filtered}
}

func (c renderComposite) Render() {
	for _, item := range c.renderers {
		item.Render()
	}
}
