package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Env() Provider {
	return New("Env", func(current *context.Runtime) any {
		if current == nil || current.Env == nil {
			return map[string]string{}
		}
		return current.Env
	})
}
