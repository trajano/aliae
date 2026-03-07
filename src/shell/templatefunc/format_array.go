package templatefunc

func FormatArray(fn any) Provider {
	return New("formatArray", fn)
}
