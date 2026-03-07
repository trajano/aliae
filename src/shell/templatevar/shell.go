package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Shell() Provider {
	return New("Shell", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.Shell
	})
}
