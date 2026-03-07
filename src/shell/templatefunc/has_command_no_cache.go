package templatefunc

func HasCommandNoCache(fn any) Provider {
	return New("hasCommandNoCache", fn)
}
