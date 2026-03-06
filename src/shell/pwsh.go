package shell

import (
	"fmt"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	PWSH       = "pwsh"
	POWERSHELL = "powershell"

	AllScope    Option = "AllScope"
	Constant    Option = "Constant"
	ReadOnly    Option = "ReadOnly"
	None        Option = "None"
	Unspecified Option = "Unspecified"

	Private Option = "Private"

	Global         Option = "Global"
	Local          Option = "Local"
	NumberedScopes Option = "Numbered scopes"
	ScriptScope    Option = "Script"
)

func IsPowerShell(shell string) bool {
	return shell == PWSH || shell == POWERSHELL
}

type pwshFormatStrategy struct{}

func (pwshFormatStrategy) FormatAlias(a *Alias) string {
	// PowerShell can't handle aliases with switches
	// unlike unix shells do so we wrap those in a function
	if a.Type == Command && strings.Contains(string(a.Value), " ") {
		a.template = `function {{ .Name }}() {
	{{ .Value }} $args
}`
		return a.render()
	}

	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `Set-Alias -Name {{ .Name }} -Value {{ formatString .Value }}{{ if .Description }} -Description {{ formatString .Description }}{{ end }}{{ if .Force }} -Force{{ end }}{{ if isPwshOption .Option }} -Option {{ formatString .Option }}{{ end }}{{ if isPwshScope .Scope }} -Scope {{ formatString .Scope }}{{ end }}` //nolint: lll
	case Function:
		a.template = `function {{ .Name }}() {
    {{ .Value }}
}`
	case Python:
		a.template = `function {{ .Name }}() {
    python -c {{ formatString .Value }} $args
}`
	case Perl:
		a.template = `function {{ .Name }}() {
    perl -e {{ formatString .Value }} $args
}`
	}

	return a.render()
}

func (pwshFormatStrategy) FormatEnv(e *Env) string {
	switch e.Type {
	case Array:
		e.template = `$env:{{ .Name }} = @({{ formatArray .Value "," }})`
	case String:
		fallthrough
	default:
		e.template = `$env:{{ .Name }} = {{ formatString .Value }}`
	}

	return e.render()
}

func (pwshFormatStrategy) FormatPath(p *Path) string {
	p.template = fmt.Sprintf(`$env:PATH = {{ formatString .Value }} + '%s' + $env:PATH`, context.PathDelimiter())
	return p.render()
}

func (pwshFormatStrategy) FormatCDPath(*CDPath) string { return "" }

func (pwshFormatStrategy) FormatLink(l *Link) string {
	l.template = `New-Item -Path {{ formatString .Name }} -ItemType SymbolicLink -Value {{ formatString .Target }} -Force | Out-Null`
	return l.render()
}

func (pwshFormatStrategy) FormatEcho(e *Echo) string {
	e.template = `$message = @"
{{ .Message }}
"@
Write-Host $message`
	return e.render()
}

func (pwshFormatStrategy) FormatCDPathCurrentDirScript() string { return "" }

func isPwshOption(option Option) bool {
	switch option { //nolint:exhaustive
	case AllScope, Constant, None, Private, ReadOnly, Unspecified:
		return true
	default:
		return false
	}
}

func isPwshScope(option Option) bool {
	switch option { //nolint:exhaustive
	case Global, Local, Private, NumberedScopes, ScriptScope:
		return true
	default:
		return false
	}
}
