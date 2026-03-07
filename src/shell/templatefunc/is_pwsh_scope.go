package templatefunc

func IsPwshScope(fn any) Provider {
	return New("isPwshScope", fn)
}
