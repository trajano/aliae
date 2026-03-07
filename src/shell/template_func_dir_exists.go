package shell

type dirExistsFuncProvider struct{}

func (dirExistsFuncProvider) Name() string {
	return "dirExists"
}

func (dirExistsFuncProvider) Func() any {
	return dirExists
}
