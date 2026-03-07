package shell

import "github.com/jandedobbeleer/aliae/src/context"

type configPathVariableProvider struct{}

func (configPathVariableProvider) Name() string {
	return "ConfigPath"
}

func (configPathVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.ConfigPath
}
