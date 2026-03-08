package templatefunc

func IsDir(fn any) Provider {
	return New("isDir", fn)
}
