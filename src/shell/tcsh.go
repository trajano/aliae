package shell

const (
	TCSH = "tcsh"
)

type tcshRenderStrategy struct{}

func (tcshRenderStrategy) RenderAlias(a *Alias) string          { return a.tcsh().render() }
func (tcshRenderStrategy) RenderEnv(e *Env) string              { return e.tcsh().render() }
func (tcshRenderStrategy) RenderPath(p *Path) string            { return p.tcsh().render() }
func (tcshRenderStrategy) RenderCDPath(p *CDPath) string        { return p.tcsh().render() }
func (tcshRenderStrategy) RenderLink(l *Link) string            { return l.tcsh().render() }
func (tcshRenderStrategy) RenderEcho(e *Echo) string            { return e.zsh().render() }
func (tcshRenderStrategy) RenderCDPathCurrentDirScript() string { return `set cdpath = ( . $cdpath );` }

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
