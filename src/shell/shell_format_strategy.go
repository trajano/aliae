package shell

import "github.com/jandedobbeleer/aliae/src/context"

// ShellFormatStrategy defines shell-specific rendering behavior for all supported section types.
//
// This interface is intentionally public so callers can depend on rendering capabilities
// without coupling to a specific shell implementation.
type ShellFormatStrategy interface {
	RenderAlias(*Alias) string
	RenderEnv(*Env) string
	RenderPath(*Path) string
	RenderCDPath(*CDPath) string
	RenderLink(*Link) string
	RenderEcho(*Echo) string
	RenderCDPathCurrentDirScript() string
}

// renderStrategy applies a Factory Method-style selection to return
// the concrete ShellFormatStrategy for the active shell runtime.
func renderStrategy() ShellFormatStrategy {
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

func (noopRenderStrategy) RenderAlias(*Alias) string            { return "" }
func (noopRenderStrategy) RenderEnv(*Env) string                { return "" }
func (noopRenderStrategy) RenderPath(*Path) string              { return "" }
func (noopRenderStrategy) RenderCDPath(*CDPath) string          { return "" }
func (noopRenderStrategy) RenderLink(*Link) string              { return "" }
func (noopRenderStrategy) RenderEcho(*Echo) string              { return "" }
func (noopRenderStrategy) RenderCDPathCurrentDirScript() string { return "" }
