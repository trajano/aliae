package shell

import "github.com/jandedobbeleer/aliae/src/context"

type envVariableProvider struct{}

func (envVariableProvider) Name() string {
	return "Env"
}

func (envVariableProvider) Value(current *context.Runtime) any {
	if current == nil || current.Env == nil {
		return map[string]string{}
	}

	return current.Env
}
