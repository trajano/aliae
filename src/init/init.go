package init

import (
	"fmt"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func Init(configPath, sh string, printOutput bool) string {
	if shell.IsPowerShell(sh) && !printOutput {
		return fmt.Sprintf("(@(& aliae init %s --config=%s --print) -join \"`n\") | Invoke-Expression", sh, configPath)
	}

	cfg.BeginInitInternalProgress(configPath)
	defer cfg.EndInitInternalProgress()

	script, err := runWithObserver(configPath, sh, nil, runOptions{
		computeVars: true,
		primeState:  true,
	})
	if err != nil {
		errorString := formatError(err)
		if sh == shell.NU {
			return createNuInit(errorString)
		}
		return errorString
	}

	if sh != shell.NU || printOutput {
		return script
	}

	return createNuInit(script)
}

func createNuInit(script string) string {
	err := shell.NuInit(script)
	if err != nil {
		return formatError(err)
	}

	return ""
}

func formatError(err error) string {
	restoreRuntime := shell.SetRuntime(context.Current)
	defer restoreRuntime()
	restoreTemplateRuntime := shell.SetTemplateRuntime(context.Current)
	defer restoreTemplateRuntime()

	message := fmt.Sprintf("aliae error:\n%s", err.Error())
	e := shell.Echo{Message: message}
	return e.Error().String()
}
