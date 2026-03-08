package templatefunc

func SetArg(fn any) Provider {
	return New("setArg", fn)
}
