package shell

const (
	FISH = "fish"
)

type fishFormatStrategy struct{}

func (fishFormatStrategy) FormatAlias(a *Alias) string {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} {{ formatString .Value }}`
	case Function:
		a.template = `function {{ .Name }}{{ if .Description }} --description {{ formatString .Description }}{{ end }}
    {{ .Value }}
end`
	case Python:
		a.template = `function {{ .Name }}{{ if .Description }} --description {{ formatString .Description }}{{ end }}
    python -c {{ formatString .Value }} $argv
end`
	case Perl:
		a.template = `function {{ .Name }}{{ if .Description }} --description {{ formatString .Description }}{{ end }}
    perl -e {{ formatString .Value }} $argv
end`
	}

	return a.render()
}

func (fishFormatStrategy) FormatEnv(e *Env) string {
	e.template = `set --global --export {{ .Name }} {{ .Value }}`
	return e.render()
}

func (fishFormatStrategy) FormatPath(p *Path) string {
	p.template = `fish_add_path {{ .Value }}`
	return p.render()
}

func (fishFormatStrategy) FormatCDPath(p *CDPath) string {
	p.template = `set -g cdpath $cdpath {{ .Value }}`
	return p.render()
}

func (fishFormatStrategy) FormatLink(l *Link) string {
	l.template = defaultLinkTemplate
	return l.render()
}

func (fishFormatStrategy) FormatEcho(e *Echo) string {
	e.template = defaultEchoTemplate
	return e.render()
}

func (fishFormatStrategy) FormatCDPathCurrentDirScript() string { return `set -g cdpath . $cdpath` }
