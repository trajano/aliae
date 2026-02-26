package shell

const (
	TCSH = "tcsh"
)

func (a *Alias) tcsh() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} '{{ .Value }}';`
	case Python:
		a.template = `alias {{ .Name }} 'python -c {{ formatString .Value }}';`
	case Perl:
		a.template = `alias {{ .Name }} 'perl -e {{ formatString .Value }}';`
	}

	return a
}

func (e *Env) tcsh() *Env {
	e.template = `setenv {{ .Name }} {{ formatString .Value }};`
	return e
}

func (l *Link) tcsh() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }};`
	l.template = template
	return l
}

func (p *Path) tcsh() *Path {
	p.template = `set path = ( {{ .Value }} $path );`
	return p
}

func (p *CDPath) tcsh() *CDPath {
	p.template = `set cdpath = ( $cdpath {{ .Value }} );`
	return p
}
