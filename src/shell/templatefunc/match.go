package templatefunc

func Match(fn any) Provider {
	return New("match", fn)
}
