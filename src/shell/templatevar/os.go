package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

func OS() Provider {
	return New("OS", func(current *context.Runtime) any {
		if current == nil {
			return ""
		}
		return current.OS
	})
}
