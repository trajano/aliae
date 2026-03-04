package shell

import "github.com/jandedobbeleer/aliae/src/context"

// ShellFormatStrategy defines shell-specific rendering behavior for all supported section types.
//
// This interface is intentionally public so callers can depend on rendering capabilities
// without coupling to a specific shell implementation.
type ShellFormatStrategy interface {
	FormatAlias(*Alias) string
	FormatEnv(*Env) string
	FormatPath(*Path) string
	FormatCDPath(*CDPath) string
	FormatLink(*Link) string
	FormatEcho(*Echo) string
	FormatCDPathCurrentDirScript() string
}

// formatStrategy applies a Factory Method-style selection to return
// the concrete ShellFormatStrategy for the active shell runtime.
func formatStrategy() ShellFormatStrategy {
	if context.Current == nil {
		return noopFormatStrategy{}
	}

	switch context.Current.Shell {
	case ZSH:
		return zshFormatStrategy{}
	case BASH:
		return bashFormatStrategy{}
	case PWSH, POWERSHELL:
		return pwshFormatStrategy{}
	case NU:
		return nuFormatStrategy{}
	case FISH:
		return fishFormatStrategy{}
	case TCSH:
		return tcshFormatStrategy{}
	case XONSH:
		return xonshFormatStrategy{}
	case CMD:
		return cmdFormatStrategy{}
	default:
		return noopFormatStrategy{}
	}
}

type noopFormatStrategy struct{}

func (noopFormatStrategy) FormatAlias(*Alias) string            { return "" }
func (noopFormatStrategy) FormatEnv(*Env) string                { return "" }
func (noopFormatStrategy) FormatPath(*Path) string              { return "" }
func (noopFormatStrategy) FormatCDPath(*CDPath) string          { return "" }
func (noopFormatStrategy) FormatLink(*Link) string              { return "" }
func (noopFormatStrategy) FormatEcho(*Echo) string              { return "" }
func (noopFormatStrategy) FormatCDPathCurrentDirScript() string { return "" }
