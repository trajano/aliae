package shell

type matchFuncProvider struct{}

func (matchFuncProvider) Name() string {
	return "match"
}

func (matchFuncProvider) Func() any {
	return match
}
