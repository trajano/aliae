package templatefunc

func HasCommand(fn any) Provider {
	return New("hasCommand", fn)
}
