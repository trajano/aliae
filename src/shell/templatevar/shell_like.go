package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func ShellLike() Provider {
	return New("ShellLike", func(current *context.Runtime) any {
		if current == nil {
			return false
		}
		return current.ShellLike
	})
}
