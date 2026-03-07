package shell

type escapeStringFuncProvider struct{}

func (escapeStringFuncProvider) Name() string {
	return "escapeString"
}

func (escapeStringFuncProvider) Func() any {
	return escapeString
}
