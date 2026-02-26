package shell

const (
	FISH = "fish"
)

func (a *Alias) fish() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} {{ formatString .Value }}`
	case Function:
		a.template = `function {{ .Name }}{{ if .Description }} --description {{ formatString .Description }}{{ end }}
    {{ .Value }}
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
