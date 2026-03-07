package shell

import "github.com/jandedobbeleer/aliae/src/context"

type wslVariableProvider struct{}

func (wslVariableProvider) Name() string {
	return "WSL"
}

func (wslVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return false
	}

	return current.WSL
}
