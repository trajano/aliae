package shell

type hasCommandFuncProvider struct{}

func (hasCommandFuncProvider) Name() string {
	return "hasCommand"
}

func (hasCommandFuncProvider) Func() any {
	return hasCommand
}
