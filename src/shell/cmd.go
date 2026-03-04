package shell

const (
	CMD = "cmd"
)

type cmdRenderStrategy struct{}

func (cmdRenderStrategy) RenderAlias(a *Alias) string          { return a.cmd().render() }
func (cmdRenderStrategy) RenderEnv(e *Env) string              { return e.cmd().render() }
func (cmdRenderStrategy) RenderPath(p *Path) string            { return p.cmd().render() }
func (cmdRenderStrategy) RenderCDPath(*CDPath) string          { return "" }
func (cmdRenderStrategy) RenderLink(l *Link) string            { return l.cmd().render() }
func (cmdRenderStrategy) RenderEcho(e *Echo) string            { return e.cmd().render() }
func (cmdRenderStrategy) RenderCDPathCurrentDirScript() string { return "" }

func (a *Alias) cmd() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = "macrofile:write(\"{{ .Name }}={{ escapeString .Value }}\", \"\\n\")"
	case Python:
		a.template = "macrofile:write(\"{{ .Name }}=python -c \\\"{{ escapeString .Value }}\\\" $*\", \"\\n\")"
	case Perl:
		a.template = "macrofile:write(\"{{ .Name }}=perl -e \\\"{{ escapeString .Value }}\\\" $*\", \"\\n\")"
	}

	return a
}

func cmdAliasPre() string {
	return `local filename  = os.tmpname()
local macrofile = io.open(filename, "w+")
`
}

func cmdAliasPost() string {
	return `
macrofile:close()
local _ = io.popen(string.format("doskey /macrofile=%s", filename)):close()
os.remove(filename)`
}

func (e *Echo) cmd() *Echo {
	e.template = `message = [[
{{ .Message }}
]]
print(message)`
	return e
}

func (e *Env) cmd() *Env {
	e.template = `os.setenv("{{ .Name }}", {{ formatString .Value }})`
	return e
}

func (l *Link) cmd() *Link {
	template := `os.execute("{{ $source := (escapeString .Name) }}mklink {{ if isDir $source }}/d{{ else }}/h{{ end }} {{ $source }} {{ escapeString .Target }} > nul 2>&1")`
	l.template = template
	return l
}

func (p *Path) cmd() *Path {
	p.template = `os.setenv("PATH", "{{ escapeString .Value }};" .. os.getenv("PATH"))`
	return p
}
