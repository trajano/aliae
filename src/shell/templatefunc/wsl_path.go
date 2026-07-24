package templatefunc

func WslPath(fn any) Provider {
	return New("wslPath", fn)
}
