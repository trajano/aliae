package shell

import "github.com/jandedobbeleer/aliae/src/context"

type templateVariableProvider interface {
	Name() string
	Value(current *context.Runtime) any
}

var templateVariableProviders = []templateVariableProvider{
	shellVariableProvider{},
	shellLikeVariableProvider{},
	wslVariableProvider{},
	homeVariableProvider{},
	osVariableProvider{},
	archVariableProvider{},
	hostnameVariableProvider{},
	configPathVariableProvider{},
	configDirVariableProvider{},
	envVariableProvider{},
	variableMapVariableProvider{},
}

func templateContext(ctx any) any {
	switch current := ctx.(type) {
	case *context.Runtime:
		return templateRuntimeVariables(current)
	case context.Runtime:
		runtimeCopy := current
		return templateRuntimeVariables(&runtimeCopy)
	default:
		return ctx
	}
}

func templateRuntimeVariables(current *context.Runtime) map[string]any {
	values := make(map[string]any, len(templateVariableProviders))

	for _, provider := range templateVariableProviders {
		values[provider.Name()] = provider.Value(current)
	}

	return values
}
