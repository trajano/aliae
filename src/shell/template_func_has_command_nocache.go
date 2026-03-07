package shell

type hasCommandNoCacheFuncProvider struct{}

func (hasCommandNoCacheFuncProvider) Name() string {
	return "hasCommandNoCache"
}

func (hasCommandNoCacheFuncProvider) Func() any {
	return hasCommandNoCache
}
