package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func Hostname() Provider {
	return New("Hostname", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.Hostname
	})
}
