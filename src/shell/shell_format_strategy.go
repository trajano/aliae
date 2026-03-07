package shell

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
	FormatSetArg(name string, oneBasedIndex int) string
	FormatProgress(state, percentage int) string
	EscapeString(value string) string
	FormatAliasScriptPrelude() string
	FormatAliasScriptPostlude() string
}

// formatStrategy applies a Factory Method-style selection to return
// the concrete ShellFormatStrategy for the active shell runtime.
func formatStrategy() ShellFormatStrategy {
	runtime := currentRuntime()
	if runtime == nil {
		return noopFormatStrategy{}
	}

	return formatStrategyForShell(runtime.Shell)
}

func formatStrategyForShell(shellName string) ShellFormatStrategy {
	switch shellName {
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
func (noopFormatStrategy) FormatSetArg(string, int) string      { return "" }
func (noopFormatStrategy) FormatProgress(int, int) string       { return "" }
func (noopFormatStrategy) EscapeString(value string) string     { return defaultEscapedString(value) }
func (noopFormatStrategy) FormatAliasScriptPrelude() string     { return "" }
func (noopFormatStrategy) FormatAliasScriptPostlude() string    { return "" }
