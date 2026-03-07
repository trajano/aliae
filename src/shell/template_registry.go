package shell

import (
	"text/template"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell/templatefunc"
	"github.com/jandedobbeleer/aliae/src/shell/templatevar"
)

var shellTemplateFuncProviders = []templatefunc.Provider{
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

var shellTemplateVariableProviders = []templatevar.Provider{
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

var shellTemplateFuncMap template.FuncMap

func init() {
	shellTemplateFuncMap = make(template.FuncMap, len(shellTemplateFuncProviders))
	for _, provider := range shellTemplateFuncProviders {
		shellTemplateFuncMap[provider.Name()] = provider.Func()
	}
}

func templateContext(ctx any) any {
	switch current := ctx.(type) {
	case *context.Runtime:
		return templatevar.Context(current, shellTemplateVariableProviders)
	case context.Runtime:
		runtimeCopy := current
		return templatevar.Context(&runtimeCopy, shellTemplateVariableProviders)
	default:
		return ctx
	}
}
