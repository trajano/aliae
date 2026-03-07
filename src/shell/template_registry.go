package shell

import (
	"os"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell/templatefunc"
	"github.com/jandedobbeleer/aliae/src/shell/templatevar"
)

func shellTemplateFuncProviders() []templatefunc.Provider {
	return []templatefunc.Provider{
		templatefunc.New("isPwshOption", isPwshOption),
		templatefunc.New("isPwshScope", isPwshScope),
		templatefunc.New("formatString", formatString),
		templatefunc.New("formatArray", formatArray),
		templatefunc.New("escapeString", escapeString),
		templatefunc.New("env", os.Getenv),
		templatefunc.New("match", match),
		templatefunc.New("hasCommand", hasCommand),
		templatefunc.New("hasCommandNoCache", hasCommandNoCache),
		templatefunc.New("fileExists", fileExists),
		templatefunc.New("dirExists", dirExists),
		templatefunc.New("isDir", isDir),
		templatefunc.New("setArg", setArg),
		templatefunc.New("progress", progress),
	}
}

func shellTemplateVariableProviders() []templatevar.Provider {
	return []templatevar.Provider{
		templatevar.New("Shell", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.Shell
		}),
		templatevar.New("ShellLike", func(current *context.Runtime) any {
			if current == nil {
				return false
			}

			return current.ShellLike
		}),
		templatevar.New("WSL", func(current *context.Runtime) any {
			if current == nil {
				return false
			}

			return current.WSL
		}),
		templatevar.New("Home", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.Home
		}),
		templatevar.New("OS", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.OS
		}),
		templatevar.New("Arch", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.Arch
		}),
		templatevar.New("Hostname", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.Hostname
		}),
		templatevar.New("ConfigPath", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.ConfigPath
		}),
		templatevar.New("ConfigDir", func(current *context.Runtime) any {
			if current == nil {
				return ""
			}

			return current.ConfigDir
		}),
		templatevar.New("Env", func(current *context.Runtime) any {
			if current == nil || current.Env == nil {
				return map[string]string{}
			}

			return current.Env
		}),
		templatevar.New("Var", func(current *context.Runtime) any {
			if current == nil || current.Var == nil {
				return map[string]any{}
			}

			return current.Var
		}),
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
