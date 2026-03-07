package shell

type progressFuncProvider struct{}

func (progressFuncProvider) Name() string {
	return "progress"
}

func (progressFuncProvider) Func() any {
	return progress
}
