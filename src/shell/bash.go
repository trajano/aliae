package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	BASH = "bash"
)

func (a *Alias) bash() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }}={{ formatString .Value }}`
	case Function:
		a.template = `{{ .Name }}() {
    {{ .Value }}
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
	template := fmt.Sprintf(`export CDPATH="{{ .Value }}%s$CDPATH"`, context.PathDelimiter())
	p.template = template
	return p
}
