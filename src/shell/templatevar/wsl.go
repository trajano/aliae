package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func WSL() Provider {
	return New("WSL", func(current *context.Runtime) any {
		if current == nil {
			return false
		}
		return current.WSL
	})
}
