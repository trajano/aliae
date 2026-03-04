package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	BASH = "bash"
)

type bashRenderStrategy struct{}

func (bashRenderStrategy) RenderAlias(a *Alias) string   { return a.bash().render() }
func (bashRenderStrategy) RenderEnv(e *Env) string       { return e.bash().render() }
func (bashRenderStrategy) RenderPath(p *Path) string     { return p.bash().render() }
func (bashRenderStrategy) RenderCDPath(p *CDPath) string { return p.bash().render() }
func (bashRenderStrategy) RenderLink(l *Link) string     { return l.bash().render() }
func (bashRenderStrategy) RenderEcho(e *Echo) string     { return e.bash().render() }
func (bashRenderStrategy) RenderCDPathCurrentDirScript() string {
	return `if [ -n "$CDPATH" ]; then export CDPATH=":$CDPATH"; else export CDPATH=":"; fi`
}

func (a *Alias) bash() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }}={{ formatString .Value }}`
	case Function:
		a.template = `{{ .Name }}() {
    {{ .Value }}
}`
	case Python:
		a.template = `{{ .Name }}() {
    python -c {{ formatString .Value }} "$@"
}`
	case Perl:
		a.template = `{{ .Name }}() {
    perl -e {{ formatString .Value }} "$@"
}`
	}

	return a
}

func (e *Echo) bash() *Echo {
	e.template = defaultEchoTemplate
	return e
}

func (e *Env) bash() *Env {
	switch e.Type {
	case Array:
		e.template = `export {{ .Name }}=({{ formatArray .Value }})`
	case String:
		fallthrough
	default:
		e.template = `export {{ .Name }}={{ formatString .Value }}`
	}

	return e
}

func (l *Link) bash() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }}`
	l.template = template
	return l
}

func (p *Path) bash() *Path {
	p.template = `export PATH="{{ .Value }}:$PATH"`
	return p
}

func (p *CDPath) bash() *CDPath {
	template := fmt.Sprintf(`export CDPATH="${CDPATH:+$CDPATH%s}{{ .Value }}"`, context.PathDelimiter())
	p.template = template
	return p
}
