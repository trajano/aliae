package shell

const (
	TCSH = "tcsh"
)

type tcshFormatStrategy struct{}

func (tcshFormatStrategy) FormatAlias(a *Alias) string {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} '{{ .Value }}';`
	case Python:
		a.template = `alias {{ .Name }} 'python -c {{ formatString .Value }}';`
	case Perl:
		a.template = `alias {{ .Name }} 'perl -e {{ formatString .Value }}';`
	}

	return a.render()
}

func (tcshFormatStrategy) FormatEnv(e *Env) string {
	e.template = `setenv {{ .Name }} {{ formatString .Value }};`
	return e.render()
}

func (tcshFormatStrategy) FormatPath(p *Path) string {
	p.template = `set path = ( {{ .Value }} $path );`
	return p.render()
}

func (tcshFormatStrategy) FormatCDPath(p *CDPath) string {
	p.template = `set cdpath = ( $cdpath {{ .Value }} );`
	return p.render()
}

func (tcshFormatStrategy) FormatLink(l *Link) string {
	l.template = `ln -sf {{ .Target }} {{ .Name }};`
	return l.render()
}

func (tcshFormatStrategy) FormatEcho(e *Echo) string {
	e.template = defaultEchoTemplate
	return e.render()
}

func (tcshFormatStrategy) FormatCDPathCurrentDirScript() string { return `set cdpath = ( . $cdpath );` }
