package shell

type isDirFuncProvider struct{}

func (isDirFuncProvider) Name() string {
	return "isDir"
}

func (isDirFuncProvider) Func() any {
	return isDir
}
