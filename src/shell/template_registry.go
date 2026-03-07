package shell

import (
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell/templatefunc"
	"github.com/jandedobbeleer/aliae/src/shell/templatevar"
)

func shellTemplateFuncProviders() []templatefunc.Provider {
	return []templatefunc.Provider{
		templatefunc.IsPwshOption(isPwshOption),
		templatefunc.IsPwshScope(isPwshScope),
		templatefunc.FormatString(formatString),
		templatefunc.FormatArray(formatArray),
		templatefunc.EscapeString(escapeString),
		templatefunc.Env(),
		templatefunc.Match(match),
		templatefunc.HasCommand(hasCommand),
		templatefunc.HasCommandNoCache(hasCommandNoCache),
		templatefunc.FileExists(fileExists),
		templatefunc.DirExists(dirExists),
		templatefunc.IsDir(isDir),
		templatefunc.SetArg(setArg),
		templatefunc.Progress(progress),
	}
}

func shellTemplateVariableProviders() []templatevar.Provider {
	return []templatevar.Provider{
		templatevar.Shell(),
		templatevar.ShellLike(),
		templatevar.WSL(),
		templatevar.Home(),
		templatevar.OS(),
		templatevar.Arch(),
		templatevar.Hostname(),
		templatevar.ConfigPath(),
		templatevar.ConfigDir(),
		templatevar.Env(),
		templatevar.Var(),
	}
}

func templateContext(ctx any) any {
	switch current := ctx.(type) {
	case *context.Runtime:
		return templatevar.Context(current, shellTemplateVariableProviders())
	case context.Runtime:
		runtimeCopy := current
		return templatevar.Context(&runtimeCopy, shellTemplateVariableProviders())
	default:
		return ctx
	}
}
