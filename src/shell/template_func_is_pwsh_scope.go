package shell

type isPwshScopeFuncProvider struct{}

func (isPwshScopeFuncProvider) Name() string {
	return "isPwshScope"
}

func (isPwshScopeFuncProvider) Func() any {
	return isPwshScope
}
