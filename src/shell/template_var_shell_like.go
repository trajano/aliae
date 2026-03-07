package shell

import "github.com/jandedobbeleer/aliae/src/context"

type shellLikeVariableProvider struct{}

func (shellLikeVariableProvider) Name() string {
	return "ShellLike"
}

func (shellLikeVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return false
	}

	return current.ShellLike
}
