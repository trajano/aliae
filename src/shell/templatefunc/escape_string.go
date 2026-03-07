package templatefunc

func EscapeString(fn any) Provider {
	return New("escapeString", fn)
}
