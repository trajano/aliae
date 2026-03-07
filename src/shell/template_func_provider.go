package shell

import "os"

type templateFuncProvider interface {
	Name() string
	Func() any
}

var templateFuncProviders = []templateFuncProvider{
	isPwshOptionFuncProvider{},
	isPwshScopeFuncProvider{},
	formatStringFuncProvider{},
	formatArrayFuncProvider{},
	escapeStringFuncProvider{},
	envFuncProvider{},
	matchFuncProvider{},
	hasCommandFuncProvider{},
	hasCommandNoCacheFuncProvider{},
	fileExistsFuncProvider{},
	dirExistsFuncProvider{},
	isDirFuncProvider{},
	setArgFuncProvider{},
	progressFuncProvider{},
}

type envFuncProvider struct{}

func (envFuncProvider) Name() string {
	return "env"
}

func (envFuncProvider) Func() any {
	return os.Getenv
}
