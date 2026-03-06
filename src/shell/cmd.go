package shell

const (
	CMD = "cmd"
)

type cmdFormatStrategy struct{}

func (cmdFormatStrategy) FormatAlias(a *Alias) string {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = "macrofile:write(\"{{ .Name }}={{ escapeString .Value }}\", \"\\n\")"
	case Python:
		a.template = "macrofile:write(\"{{ .Name }}=python -c \\\"{{ escapeString .Value }}\\\" $*\", \"\\n\")"
	case Perl:
		a.template = "macrofile:write(\"{{ .Name }}=perl -e \\\"{{ escapeString .Value }}\\\" $*\", \"\\n\")"
	}

	return a.render()
}

func (cmdFormatStrategy) FormatEnv(e *Env) string {
	e.template = `os.setenv("{{ .Name }}", {{ formatString .Value }})`
	return e.render()
}

func (cmdFormatStrategy) FormatPath(p *Path) string {
	p.template = `os.setenv("PATH", "{{ escapeString .Value }};" .. os.getenv("PATH"))`
	return p.render()
}

func (cmdFormatStrategy) FormatCDPath(*CDPath) string { return "" }

func (cmdFormatStrategy) FormatLink(l *Link) string {
	l.template = `os.execute("{{ $source := (escapeString .Name) }}mklink {{ if isDir $source }}/d{{ else }}/h{{ end }} {{ $source }} {{ escapeString .Target }} > nul 2>&1")`
	return l.render()
}

func (cmdFormatStrategy) FormatEcho(e *Echo) string {
	e.template = `message = [[
{{ .Message }}
]]
print(message)`
	return e.render()
}

func (cmdFormatStrategy) FormatCDPathCurrentDirScript() string { return "" }

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
