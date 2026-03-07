package shell

import "github.com/jandedobbeleer/aliae/src/context"

type variableMapVariableProvider struct{}

func (variableMapVariableProvider) Name() string {
	return "Var"
}

func (variableMapVariableProvider) Value(current *context.Runtime) any {
	if current == nil || current.Var == nil {
		return map[string]string{}
	}

	return current.Var
}
