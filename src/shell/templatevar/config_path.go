package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func ConfigPath() Provider {
	return New("ConfigPath", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.ConfigPath
	})
}
