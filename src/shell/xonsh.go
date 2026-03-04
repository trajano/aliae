package shell

import (
	"fmt"
	"strings"
)

const (
	XONSH = "xonsh"
)

type xonshRenderStrategy struct{}

func (xonshRenderStrategy) RenderAlias(a *Alias) string          { return a.xonsh().render() }
func (xonshRenderStrategy) RenderEnv(e *Env) string              { return e.xonsh().render() }
func (xonshRenderStrategy) RenderPath(p *Path) string            { return p.xonsh().render() }
func (xonshRenderStrategy) RenderCDPath(p *CDPath) string        { return p.xonsh().render() }
func (xonshRenderStrategy) RenderLink(l *Link) string            { return l.zsh().render() }
func (xonshRenderStrategy) RenderEcho(e *Echo) string            { return e.xonsh().render() }
func (xonshRenderStrategy) RenderCDPathCurrentDirScript() string { return `$CDPATH = ["."] + $CDPATH` }

func (a *Alias) xonsh() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `aliases['{{ .Name }}'] = '{{ .Value }}'`
	case Function:
		// some xonsh aliases are not valid python function names
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s():
    {{ .Value }}`, funcName)
		a.template = template
	case Python:
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s(args):
    import subprocess
    subprocess.run(["python", "-c", {{ formatString .Value }}, *args], check=False)`, funcName)
		a.template = template
	case Perl:
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s(args):
    import subprocess
    subprocess.run(["perl", "-e", {{ formatString .Value }}, *args], check=False)`, funcName)
		a.template = template
	}

	return a
}

func (e *Echo) xonsh() *Echo {
	e.template = `message = """{{ .Message }}"""
print(message)`
	return e
}

func (e *Env) xonsh() *Env {
	switch e.Type {
	case Array:
		e.template = `${{ .Name }} = [{{ formatArray .Value "," }}]`
	case String:
		fallthrough
	default:
		e.template = `${{ .Name }} = {{ formatString .Value }}`
	}

	return e
}

func (p *Path) xonsh() *Path {
	p.template = `$PATH.add('{{ .Value }}', True, False)`
	return p
}

func (p *CDPath) xonsh() *CDPath {
	p.template = `$CDPATH = $CDPATH + [{{ formatString .Value }}]`
	return p
}
