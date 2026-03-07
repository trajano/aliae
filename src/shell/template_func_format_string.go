package shell

type formatStringFuncProvider struct{}

func (formatStringFuncProvider) Name() string {
	return "formatString"
}

func (formatStringFuncProvider) Func() any {
	return formatString
}
