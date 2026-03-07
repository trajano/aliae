package shell

import "github.com/jandedobbeleer/aliae/src/context"

type shellVariableProvider struct{}

func (shellVariableProvider) Name() string {
	return "Shell"
}

func (shellVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.Shell
}
