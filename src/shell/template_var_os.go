package shell

import "github.com/jandedobbeleer/aliae/src/context"

type osVariableProvider struct{}

func (osVariableProvider) Name() string {
	return "OS"
}

func (osVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.OS
}
