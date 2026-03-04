package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	ZSH = "zsh"
)

type zshFormatStrategy struct{}

func (zshFormatStrategy) FormatAlias(a *Alias) string {
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

func (zshFormatStrategy) FormatEnv(e *Env) string {
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

func (zshFormatStrategy) FormatPath(p *Path) string {
	p.template = fmt.Sprintf(`export PATH="{{ .Value }}%s$PATH"`, context.PathDelimiter())
	return p.render()
}

func (zshFormatStrategy) FormatCDPath(p *CDPath) string {
	p.template = `cdpath=( $cdpath {{ .Value }} )`
	return p.render()
}

func (zshFormatStrategy) FormatLink(l *Link) string {
	l.template = `ln -sf {{ .Target }} {{ .Name }}`
	return l.render()
}

func (zshFormatStrategy) FormatEcho(e *Echo) string {
	e.template = defaultEchoTemplate
	return e.render()
}

func (zshFormatStrategy) FormatCDPathCurrentDirScript() string { return `cdpath=( . $cdpath )` }
