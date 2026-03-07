package shell

import "github.com/jandedobbeleer/aliae/src/context"

type configDirVariableProvider struct{}

func (configDirVariableProvider) Name() string {
	return "ConfigDir"
}

func (configDirVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.ConfigDir
}
