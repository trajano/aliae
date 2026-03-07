package templatefunc

func Progress(fn any) Provider {
	return New("progress", fn)
}
