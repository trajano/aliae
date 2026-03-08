package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Arch() Provider {
	return New("Arch", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.Arch
	})
}
