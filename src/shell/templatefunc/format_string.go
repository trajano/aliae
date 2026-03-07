package templatefunc

func FormatString(fn any) Provider {
	return New("formatString", fn)
}
