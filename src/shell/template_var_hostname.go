package shell

import "github.com/jandedobbeleer/aliae/src/context"

type hostnameVariableProvider struct{}

func (hostnameVariableProvider) Name() string {
	return "Hostname"
}

func (hostnameVariableProvider) Value(current *context.Runtime) any {
	if current == nil {
		return ""
	}

	return current.Hostname
}
