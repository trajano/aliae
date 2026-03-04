package shell

import "github.com/jandedobbeleer/aliae/src/context"

type shellRenderStrategy interface {
	renderAlias(*Alias) string
	renderEnv(*Env) string
	renderPath(*Path) string
	renderCDPath(*CDPath) string
	renderLink(*Link) string
	renderEcho(*Echo) string
	renderCDPathCurrentDirScript() string
}

func renderStrategy() shellRenderStrategy {
	if context.Current == nil {
		return noopRenderStrategy{}
	}

	switch context.Current.Shell {
	case ZSH:
		return zshRenderStrategy{}
	case BASH:
		return bashRenderStrategy{}
	case PWSH, POWERSHELL:
		return pwshRenderStrategy{}
	case NU:
		return nuRenderStrategy{}
	case FISH:
		return fishRenderStrategy{}
	case TCSH:
		return tcshRenderStrategy{}
	case XONSH:
		return xonshRenderStrategy{}
	case CMD:
		return cmdRenderStrategy{}
	default:
		return noopRenderStrategy{}
	}
}

type noopRenderStrategy struct{}

func (noopRenderStrategy) renderAlias(*Alias) string            { return "" }
func (noopRenderStrategy) renderEnv(*Env) string                { return "" }
func (noopRenderStrategy) renderPath(*Path) string              { return "" }
func (noopRenderStrategy) renderCDPath(*CDPath) string          { return "" }
func (noopRenderStrategy) renderLink(*Link) string              { return "" }
func (noopRenderStrategy) renderEcho(*Echo) string              { return "" }
func (noopRenderStrategy) renderCDPathCurrentDirScript() string { return "" }

type zshRenderStrategy struct{}

func (zshRenderStrategy) renderAlias(a *Alias) string          { return a.zsh().render() }
func (zshRenderStrategy) renderEnv(e *Env) string              { return e.zsh().render() }
func (zshRenderStrategy) renderPath(p *Path) string            { return p.zsh().render() }
func (zshRenderStrategy) renderCDPath(p *CDPath) string        { return p.zsh().render() }
func (zshRenderStrategy) renderLink(l *Link) string            { return l.zsh().render() }
func (zshRenderStrategy) renderEcho(e *Echo) string            { return e.zsh().render() }
func (zshRenderStrategy) renderCDPathCurrentDirScript() string { return `cdpath=( . $cdpath )` }

type bashRenderStrategy struct{}

func (bashRenderStrategy) renderAlias(a *Alias) string   { return a.bash().render() }
func (bashRenderStrategy) renderEnv(e *Env) string       { return e.bash().render() }
func (bashRenderStrategy) renderPath(p *Path) string     { return p.bash().render() }
func (bashRenderStrategy) renderCDPath(p *CDPath) string { return p.bash().render() }
func (bashRenderStrategy) renderLink(l *Link) string     { return l.bash().render() }
func (bashRenderStrategy) renderEcho(e *Echo) string     { return e.bash().render() }
func (bashRenderStrategy) renderCDPathCurrentDirScript() string {
	return `if [ -n "$CDPATH" ]; then export CDPATH=":$CDPATH"; else export CDPATH=":"; fi`
}

type pwshRenderStrategy struct{}

func (pwshRenderStrategy) renderAlias(a *Alias) string          { return a.pwsh().render() }
func (pwshRenderStrategy) renderEnv(e *Env) string              { return e.pwsh().render() }
func (pwshRenderStrategy) renderPath(p *Path) string            { return p.pwsh().render() }
func (pwshRenderStrategy) renderCDPath(*CDPath) string          { return "" }
func (pwshRenderStrategy) renderLink(l *Link) string            { return l.pwsh().render() }
func (pwshRenderStrategy) renderEcho(e *Echo) string            { return e.pwsh().render() }
func (pwshRenderStrategy) renderCDPathCurrentDirScript() string { return "" }

type nuRenderStrategy struct{}

func (nuRenderStrategy) renderAlias(a *Alias) string          { return a.nu().render() }
func (nuRenderStrategy) renderEnv(e *Env) string              { return e.nu().render() }
func (nuRenderStrategy) renderPath(p *Path) string            { return p.nu().render() }
func (nuRenderStrategy) renderCDPath(*CDPath) string          { return "" }
func (nuRenderStrategy) renderLink(l *Link) string            { return l.nu().render() }
func (nuRenderStrategy) renderEcho(e *Echo) string            { return e.nu().render() }
func (nuRenderStrategy) renderCDPathCurrentDirScript() string { return "" }

type fishRenderStrategy struct{}

func (fishRenderStrategy) renderAlias(a *Alias) string          { return a.fish().render() }
func (fishRenderStrategy) renderEnv(e *Env) string              { return e.fish().render() }
func (fishRenderStrategy) renderPath(p *Path) string            { return p.fish().render() }
func (fishRenderStrategy) renderCDPath(p *CDPath) string        { return p.fish().render() }
func (fishRenderStrategy) renderLink(l *Link) string            { return l.zsh().render() }
func (fishRenderStrategy) renderEcho(e *Echo) string            { return e.zsh().render() }
func (fishRenderStrategy) renderCDPathCurrentDirScript() string { return `set -g cdpath . $cdpath` }

type tcshRenderStrategy struct{}

func (tcshRenderStrategy) renderAlias(a *Alias) string          { return a.tcsh().render() }
func (tcshRenderStrategy) renderEnv(e *Env) string              { return e.tcsh().render() }
func (tcshRenderStrategy) renderPath(p *Path) string            { return p.tcsh().render() }
func (tcshRenderStrategy) renderCDPath(p *CDPath) string        { return p.tcsh().render() }
func (tcshRenderStrategy) renderLink(l *Link) string            { return l.tcsh().render() }
func (tcshRenderStrategy) renderEcho(e *Echo) string            { return e.zsh().render() }
func (tcshRenderStrategy) renderCDPathCurrentDirScript() string { return `set cdpath = ( . $cdpath );` }

type xonshRenderStrategy struct{}

func (xonshRenderStrategy) renderAlias(a *Alias) string          { return a.xonsh().render() }
func (xonshRenderStrategy) renderEnv(e *Env) string              { return e.xonsh().render() }
func (xonshRenderStrategy) renderPath(p *Path) string            { return p.xonsh().render() }
func (xonshRenderStrategy) renderCDPath(p *CDPath) string        { return p.xonsh().render() }
func (xonshRenderStrategy) renderLink(l *Link) string            { return l.zsh().render() }
func (xonshRenderStrategy) renderEcho(e *Echo) string            { return e.xonsh().render() }
func (xonshRenderStrategy) renderCDPathCurrentDirScript() string { return `$CDPATH = ["."] + $CDPATH` }

type cmdRenderStrategy struct{}

func (cmdRenderStrategy) renderAlias(a *Alias) string          { return a.cmd().render() }
func (cmdRenderStrategy) renderEnv(e *Env) string              { return e.cmd().render() }
func (cmdRenderStrategy) renderPath(p *Path) string            { return p.cmd().render() }
func (cmdRenderStrategy) renderCDPath(*CDPath) string          { return "" }
func (cmdRenderStrategy) renderLink(l *Link) string            { return l.cmd().render() }
func (cmdRenderStrategy) renderEcho(e *Echo) string            { return e.cmd().render() }
func (cmdRenderStrategy) renderCDPathCurrentDirScript() string { return "" }
