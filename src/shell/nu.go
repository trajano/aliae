package shell

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	NU              = "nu"
	NuEnvBlockStart = "export-env {\n"
	NuEnvBlockEnd   = "\n}"
)

type nuFormatStrategy struct{}

func (nuFormatStrategy) FormatAlias(a *Alias) string {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} = {{ .Value }}`
	case Function:
		a.template = `{{ if .Description }}# {{ .Description }}
{{ end }}def {{ .Name }} [] {
    {{ .Value }}
}`
	case Python:
		a.template = `{{ if .Description }}# {{ .Description }}
{{ end }}def {{ .Name }} [...args] {
    python -c {{ formatString .Value }} ...$args
}`
	case Perl:
		a.template = `{{ if .Description }}# {{ .Description }}
{{ end }}def {{ .Name }} [...args] {
    perl -e {{ formatString .Value }} ...$args
}`
	}

	return a.render()
}

func (nuFormatStrategy) FormatEnv(e *Env) string {
	switch e.Type {
	case Array:
		e.template = `    $env.{{ .Name }} = [{{ formatArray .Value }}]`
	case String:
		fallthrough
	default:
		e.template = `    $env.{{ .Name }} = {{ formatString .Value }}`
	}

	return e.render()
}

func (nuFormatStrategy) FormatPath(p *Path) string {
	template := `$env.%s = ($env.%s | prepend {{ formatString .Value }})`
	pathName := "PATH"
	runtime := currentRuntime()

	if runtime != nil && runtime.OS == context.WINDOWS {
		pathName = "Path"
	}

	p.template = fmt.Sprintf(template, pathName, pathName)
	return p.render()
}

func (nuFormatStrategy) FormatCDPath(*CDPath) string { return "" }

func (nuFormatStrategy) FormatLink(l *Link) string {
	l.template = `ln -sf {{ .Target }} {{ .Name }} out+err>| ignore`
	runtime := currentRuntime()
	if runtime != nil && runtime.OS == context.WINDOWS {
		l.template = `{{ $source := (escapeString .Name) }}mklink {{ if isDir $source }}/d{{ else }}/h{{ end }} {{ $source }} {{ escapeString .Target }} out+err>| ignore`
	}

	return l.render()
}

func (nuFormatStrategy) FormatEcho(e *Echo) string {
	e.template = defaultEchoTemplate
	return e.render()
}

func (nuFormatStrategy) FormatCDPathCurrentDirScript() string { return "" }

func NuInit(script string) error {
	initPath := filepath.Join(context.Home(), ".aliae.nu")

	f, err := os.OpenFile(initPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = f.WriteString(script)
	if err != nil {
		return err
	}

	return f.Close()
}
