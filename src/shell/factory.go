package shell

import "github.com/jandedobbeleer/aliae/src/context"

type shellFactory interface {
	configureAlias(*Alias) *Alias
	configureEnv(*Env) *Env
	configurePath(*Path) *Path
	configureCDPath(*CDPath) *CDPath
	configureLink(*Link) *Link
	configureEcho(*Echo) *Echo
}

func newShellFactory() shellFactory {
	if context.Current == nil {
		return noopShellFactory{}
	}

	switch context.Current.Shell {
	case ZSH:
		return zshShellFactory{}
	case BASH:
		return bashShellFactory{}
	case PWSH, POWERSHELL:
		return pwshShellFactory{}
	case NU:
		return nuShellFactory{}
	case FISH:
		return fishShellFactory{}
	case TCSH:
		return tcshShellFactory{}
	case XONSH:
		return xonshShellFactory{}
	case CMD:
		return cmdShellFactory{}
	default:
		return noopShellFactory{}
	}
}

type noopShellFactory struct{}

func (noopShellFactory) configureAlias(*Alias) *Alias    { return nil }
func (noopShellFactory) configureEnv(*Env) *Env          { return nil }
func (noopShellFactory) configurePath(*Path) *Path       { return nil }
func (noopShellFactory) configureCDPath(*CDPath) *CDPath { return nil }
func (noopShellFactory) configureLink(*Link) *Link       { return nil }
func (noopShellFactory) configureEcho(*Echo) *Echo       { return nil }

type zshShellFactory struct{}

func (zshShellFactory) configureAlias(a *Alias) *Alias    { return a.zsh() }
func (zshShellFactory) configureEnv(e *Env) *Env          { return e.zsh() }
func (zshShellFactory) configurePath(p *Path) *Path       { return p.zsh() }
func (zshShellFactory) configureCDPath(p *CDPath) *CDPath { return p.zsh() }
func (zshShellFactory) configureLink(l *Link) *Link       { return l.zsh() }
func (zshShellFactory) configureEcho(e *Echo) *Echo       { return e.zsh() }

type bashShellFactory struct{}

func (bashShellFactory) configureAlias(a *Alias) *Alias    { return a.bash() }
func (bashShellFactory) configureEnv(e *Env) *Env          { return e.bash() }
func (bashShellFactory) configurePath(p *Path) *Path       { return p.bash() }
func (bashShellFactory) configureCDPath(p *CDPath) *CDPath { return p.bash() }
func (bashShellFactory) configureLink(l *Link) *Link       { return l.bash() }
func (bashShellFactory) configureEcho(e *Echo) *Echo       { return e.bash() }

type pwshShellFactory struct{}

func (pwshShellFactory) configureAlias(a *Alias) *Alias  { return a.pwsh() }
func (pwshShellFactory) configureEnv(e *Env) *Env        { return e.pwsh() }
func (pwshShellFactory) configurePath(p *Path) *Path     { return p.pwsh() }
func (pwshShellFactory) configureCDPath(*CDPath) *CDPath { return nil }
func (pwshShellFactory) configureLink(l *Link) *Link     { return l.pwsh() }
func (pwshShellFactory) configureEcho(e *Echo) *Echo     { return e.pwsh() }

type nuShellFactory struct{}

func (nuShellFactory) configureAlias(a *Alias) *Alias  { return a.nu() }
func (nuShellFactory) configureEnv(e *Env) *Env        { return e.nu() }
func (nuShellFactory) configurePath(p *Path) *Path     { return p.nu() }
func (nuShellFactory) configureCDPath(*CDPath) *CDPath { return nil }
func (nuShellFactory) configureLink(l *Link) *Link     { return l.nu() }
func (nuShellFactory) configureEcho(e *Echo) *Echo     { return e.nu() }

type fishShellFactory struct{}

func (fishShellFactory) configureAlias(a *Alias) *Alias    { return a.fish() }
func (fishShellFactory) configureEnv(e *Env) *Env          { return e.fish() }
func (fishShellFactory) configurePath(p *Path) *Path       { return p.fish() }
func (fishShellFactory) configureCDPath(p *CDPath) *CDPath { return p.fish() }
func (fishShellFactory) configureLink(l *Link) *Link       { return l.zsh() }
func (fishShellFactory) configureEcho(e *Echo) *Echo       { return e.zsh() }

type tcshShellFactory struct{}

func (tcshShellFactory) configureAlias(a *Alias) *Alias    { return a.tcsh() }
func (tcshShellFactory) configureEnv(e *Env) *Env          { return e.tcsh() }
func (tcshShellFactory) configurePath(p *Path) *Path       { return p.tcsh() }
func (tcshShellFactory) configureCDPath(p *CDPath) *CDPath { return p.tcsh() }
func (tcshShellFactory) configureLink(l *Link) *Link       { return l.tcsh() }
func (tcshShellFactory) configureEcho(e *Echo) *Echo       { return e.zsh() }

type xonshShellFactory struct{}

func (xonshShellFactory) configureAlias(a *Alias) *Alias    { return a.xonsh() }
func (xonshShellFactory) configureEnv(e *Env) *Env          { return e.xonsh() }
func (xonshShellFactory) configurePath(p *Path) *Path       { return p.xonsh() }
func (xonshShellFactory) configureCDPath(p *CDPath) *CDPath { return p.xonsh() }
func (xonshShellFactory) configureLink(l *Link) *Link       { return l.zsh() }
func (xonshShellFactory) configureEcho(e *Echo) *Echo       { return e.xonsh() }

type cmdShellFactory struct{}

func (cmdShellFactory) configureAlias(a *Alias) *Alias  { return a.cmd() }
func (cmdShellFactory) configureEnv(e *Env) *Env        { return e.cmd() }
func (cmdShellFactory) configurePath(p *Path) *Path     { return p.cmd() }
func (cmdShellFactory) configureCDPath(*CDPath) *CDPath { return nil }
func (cmdShellFactory) configureLink(l *Link) *Link     { return l.cmd() }
func (cmdShellFactory) configureEcho(e *Echo) *Echo     { return e.cmd() }
