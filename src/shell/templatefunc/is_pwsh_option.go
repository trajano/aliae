package templatefunc

func IsPwshOption(fn any) Provider {
	return New("isPwshOption", fn)
}
