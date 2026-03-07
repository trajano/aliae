package templatefunc

func FileExists(fn any) Provider {
	return New("fileExists", fn)
}
