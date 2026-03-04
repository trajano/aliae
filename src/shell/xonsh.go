package shell

import (
	"fmt"
	"strings"
)

const (
	XONSH = "xonsh"
)

type xonshFormatStrategy struct{}

func (xonshFormatStrategy) FormatAlias(a *Alias) string {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `aliases['{{ .Name }}'] = '{{ .Value }}'`
	case Function:
		// some xonsh aliases are not valid python function names
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		a.template = fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s():
    {{ .Value }}`, funcName)
	case Python:
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		a.template = fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s(args):
    import subprocess
    subprocess.run(["python", "-c", {{ formatString .Value }}, *args], check=False)`, funcName)
	case Perl:
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		a.template = fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s(args):
    import subprocess
    subprocess.run(["perl", "-e", {{ formatString .Value }}, *args], check=False)`, funcName)
	}

	return a.render()
}

func (xonshFormatStrategy) FormatEnv(e *Env) string {
	switch e.Type {
	case Array:
		e.template = `${{ .Name }} = [{{ formatArray .Value "," }}]`
	case String:
		fallthrough
	default:
		e.template = `${{ .Name }} = {{ formatString .Value }}`
	}

	return e.render()
}

func (xonshFormatStrategy) FormatPath(p *Path) string {
	p.template = `$PATH.add('{{ .Value }}', True, False)`
	return p.render()
}

func (xonshFormatStrategy) FormatCDPath(p *CDPath) string {
	p.template = `$CDPATH = $CDPATH + [{{ formatString .Value }}]`
	return p.render()
}

func (xonshFormatStrategy) FormatLink(l *Link) string {
	l.template = `ln -sf {{ .Target }} {{ .Name }}`
	return l.render()
}

func (xonshFormatStrategy) FormatEcho(e *Echo) string {
	e.template = `message = """{{ .Message }}"""
print(message)`
	return e.render()
}

func (xonshFormatStrategy) FormatCDPathCurrentDirScript() string { return `$CDPATH = ["."] + $CDPATH` }
