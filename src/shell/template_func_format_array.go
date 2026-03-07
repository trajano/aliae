package shell

type formatArrayFuncProvider struct{}

func (formatArrayFuncProvider) Name() string {
	return "formatArray"
}

func (formatArrayFuncProvider) Func() any {
	return formatArray
}
