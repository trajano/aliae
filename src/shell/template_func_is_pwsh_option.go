package shell

type isPwshOptionFuncProvider struct{}

func (isPwshOptionFuncProvider) Name() string {
	return "isPwshOption"
}

func (isPwshOptionFuncProvider) Func() any {
	return isPwshOption
}
