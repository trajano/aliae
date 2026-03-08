package templatefunc

func DirExists(fn any) Provider {
	return New("dirExists", fn)
}
