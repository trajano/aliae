package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func ConfigDir() Provider {
	return New("ConfigDir", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.ConfigDir
	})
}
