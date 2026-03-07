package templatefunc

type Provider interface {
	Name() string
	Func() any
}

type provider struct {
	name string
	fn   any
}

func New(name string, fn any) Provider {
	return provider{
		name: name,
		fn:   fn,
	}
}

func (p provider) Name() string {
	return p.name
}

func (p provider) Func() any {
	return p.fn
}
