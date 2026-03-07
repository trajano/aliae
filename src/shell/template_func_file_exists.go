package shell

type fileExistsFuncProvider struct{}

func (fileExistsFuncProvider) Name() string {
	return "fileExists"
}

func (fileExistsFuncProvider) Func() any {
	return fileExists
}
