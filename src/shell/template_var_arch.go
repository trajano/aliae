package shell

import "github.com/jandedobbeleer/aliae/src/context"

type archVariableProvider struct{}

func (archVariableProvider) Name() string {
	return "Arch"
}

func (archVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.Arch
}
