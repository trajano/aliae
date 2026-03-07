package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Var() Provider {
	return New("Var", func(current *context.Runtime) any {
		if current == nil || current.Var == nil {
			return map[string]any{}
		}
		return current.Var
	})
}
