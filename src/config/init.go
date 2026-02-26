package config

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func Init(configPath, sh string, printOutput bool) string {
	if shell.IsPowerShell(sh) && !printOutput {
		return fmt.Sprintf("(@(& aliae init %s --config=%s --print) -join \"`n\") | Invoke-Expression", sh, configPath)
	}

	beginInitInternalProgress(configPath)
	defer endInitInternalProgress()

	context.Init(sh)

	aliae, err := LoadConfig(configPath)
	if err != nil {
		errorString := formatError(err)
		if sh == shell.NU {
			return createNuInit(errorString)
		}
		return errorString
	}
	markInternalProgressConfigValidated()

	shell.StartAutoProgress(aliae.autoProgressConfig())

	aliae.Envs.Render()
	aliae.Paths.Render()
	aliae.CDPaths.Render()
	aliae.Aliae.Render()
	aliae.Links.Render()
	aliae.Scripts.Render()
	shell.EndAutoProgress()
	markInternalProgressStatPhaseComplete()

	markInternalProgressOutputFormulated()
	script := shell.DotFile.String()
	markInternalProgressReadyToOutput()

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
	message := fmt.Sprintf("aliae error:\n%s", err.Error())
	e := shell.Echo{Message: message}
	return e.Error().String()
}
