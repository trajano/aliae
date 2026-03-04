package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	ZSH = "zsh"
)

type zshRenderStrategy struct{}

func (zshRenderStrategy) RenderAlias(a *Alias) string          { return a.zsh().render() }
func (zshRenderStrategy) RenderEnv(e *Env) string              { return e.zsh().render() }
func (zshRenderStrategy) RenderPath(p *Path) string            { return p.zsh().render() }
func (zshRenderStrategy) RenderCDPath(p *CDPath) string        { return p.zsh().render() }
func (zshRenderStrategy) RenderLink(l *Link) string            { return l.zsh().render() }
func (zshRenderStrategy) RenderEcho(e *Echo) string            { return e.zsh().render() }
func (zshRenderStrategy) RenderCDPathCurrentDirScript() string { return `cdpath=( . $cdpath )` }

func (a *Alias) zsh() *Alias {
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

func (e *Echo) zsh() *Echo {
	e.template = defaultEchoTemplate
	return e
}

func (e *Env) zsh() *Env {
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

func (l *Link) zsh() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }}`
	l.template = template
	return l
}

func (p *Path) zsh() *Path {
	template := fmt.Sprintf(`export PATH="{{ .Value }}%s$PATH"`, context.PathDelimiter())
	p.template = template
	return p
}

func (p *CDPath) zsh() *CDPath {
	p.template = `cdpath=( $cdpath {{ .Value }} )`
	return p
}
