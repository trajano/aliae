package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Home() Provider {
	return New("Home", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.Home
	})
}
