package shell

const (
	FISH = "fish"
)

type fishRenderStrategy struct{}

func (fishRenderStrategy) RenderAlias(a *Alias) string          { return a.fish().render() }
func (fishRenderStrategy) RenderEnv(e *Env) string              { return e.fish().render() }
func (fishRenderStrategy) RenderPath(p *Path) string            { return p.fish().render() }
func (fishRenderStrategy) RenderCDPath(p *CDPath) string        { return p.fish().render() }
func (fishRenderStrategy) RenderLink(l *Link) string            { return l.zsh().render() }
func (fishRenderStrategy) RenderEcho(e *Echo) string            { return e.zsh().render() }
func (fishRenderStrategy) RenderCDPathCurrentDirScript() string { return `set -g cdpath . $cdpath` }

func (a *Alias) fish() *Alias {
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

	return a
}

func (e *Env) fish() *Env {
	e.template = `set --global --export {{ .Name }} {{ .Value }}`
	return e
}

func (e *Path) fish() *Path {
	e.template = `fish_add_path {{ .Value }}`
	return e
}

func (e *CDPath) fish() *CDPath {
	e.template = `set -g cdpath $cdpath {{ .Value }}`
	return e
}
