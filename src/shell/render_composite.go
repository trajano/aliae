package shell

// Renderer is the shared contract for shell section renderers.
type Renderer interface {
	Render()
}

// RenderComposite composes multiple renderers and executes them in order.
type RenderComposite struct {
	renderers []Renderer
}

func NewRenderComposite(renderers ...Renderer) RenderComposite {
	filtered := make([]Renderer, 0, len(renderers))
	for _, renderer := range renderers {
		if renderer == nil {
			continue
		}

		filtered = append(filtered, renderer)
	}

	return RenderComposite{renderers: filtered}
}

func (c RenderComposite) Render() {
	for _, renderer := range c.renderers {
		renderer.Render()
	}
}
