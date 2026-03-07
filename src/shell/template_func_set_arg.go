package shell

type setArgFuncProvider struct{}

func (setArgFuncProvider) Name() string {
	return "setArg"
}

func (setArgFuncProvider) Func() any {
	return setArg
}
