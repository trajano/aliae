package shell

import "github.com/jandedobbeleer/aliae/src/context"

type homeVariableProvider struct{}

func (homeVariableProvider) Name() string {
	return "Home"
}

func (homeVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.Home
}
