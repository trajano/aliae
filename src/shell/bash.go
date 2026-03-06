package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	BASH = "bash"
)

type bashFormatStrategy struct{}

func (bashFormatStrategy) FormatAlias(a *Alias) string {
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

	return a.render()
}

func (bashFormatStrategy) FormatEnv(e *Env) string {
	switch e.Type {
	case Array:
		e.template = `export {{ .Name }}=({{ formatArray .Value }})`
	case String:
		fallthrough
	default:
		e.template = `export {{ .Name }}={{ formatString .Value }}`
	}

	return e.render()
}

func (bashFormatStrategy) FormatPath(p *Path) string {
	p.template = `export PATH="{{ .Value }}:$PATH"`
	return p.render()
}

func (bashFormatStrategy) FormatCDPath(p *CDPath) string {
	p.template = fmt.Sprintf(`export CDPATH="${CDPATH:+$CDPATH%s}{{ .Value }}"`, context.PathDelimiter())
	return p.render()
}

func (bashFormatStrategy) FormatLink(l *Link) string {
	l.template = defaultLinkTemplate
	return l.render()
}

func (bashFormatStrategy) FormatEcho(e *Echo) string {
	e.template = defaultEchoTemplate
	return e.render()
}

func (bashFormatStrategy) FormatCDPathCurrentDirScript() string {
	return `if [ -n "$CDPATH" ]; then export CDPATH=":$CDPATH"; else export CDPATH=":"; fi`
}
