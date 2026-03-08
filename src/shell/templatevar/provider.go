package templatevar

import "github.com/jandedobbeleer/aliae/src/context"

type Provider interface {
	Name() string
	Value(current *context.Runtime) any
}

type provider struct {
	value func(current *context.Runtime) any
	name  string
}

func New(name string, value func(current *context.Runtime) any) Provider {
	return provider{
		name:  name,
		value: value,
	}
}

func (p provider) Name() string {
	return p.name
}

func (p provider) Value(current *context.Runtime) any {
	return p.value(current)
}

func Context(current *context.Runtime, providers []Provider) map[string]any {
	values := make(map[string]any, len(providers))
	for _, provider := range providers {
		values[provider.Name()] = provider.Value(current)
	}
	return values
}
