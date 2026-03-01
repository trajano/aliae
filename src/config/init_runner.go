package config

import (
	"strings"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type initRunOptions struct {
	computeVars bool
	primeState  bool
}

func runInitWithObserver(configPath, sh string, observer InitObserver, options initRunOptions) (string, error) {
	shell.DotFile.Reset()
	defer shell.DotFile.Reset()

	var aliae *Aliae
	if err := runInitPhase(observer, InitPhaseContextInit, func() error {
		context.Init(sh)
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseLoadConfig, func() error {
		var err error
		if options.computeVars {
			aliae, err = LoadConfig(configPath)
		} else {
			aliae, err = LoadConfigWithoutVars(configPath)
		}
		return err
	}); err != nil {
		return "", err
	}

	if options.computeVars {
		markInternalProgressConfigValidated()
	} else {
		if err := runInitPhase(observer, InitPhaseEvaluateVars, func() error {
			return aliae.ComputeVars()
		}); err != nil {
			return "", err
		}
	}

	if options.primeState {
		stateChecks := aliae.Scripts.PrimeState(time.Now())
		markInternalProgressStateChecksComplete(stateChecks)
	} else {
		if err := runInitPhase(observer, InitPhaseStatePrime, func() error {
			_ = aliae.Scripts.PrimeState(time.Now())
			return nil
		}); err != nil {
			return "", err
		}
	}

	if err := runInitPhase(observer, InitPhaseAutoProgressOn, func() error {
		shell.StartAutoProgress(aliae.autoProgressConfig())
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderEnv, func() error {
		aliae.Envs.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderPath, func() error {
		aliae.Paths.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderCDPath, func() error {
		aliae.CDPaths.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderAlias, func() error {
		aliae.Aliae.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderLink, func() error {
		aliae.Links.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseRenderScript, func() error {
		aliae.Scripts.Render()
		return nil
	}); err != nil {
		return "", err
	}

	if err := runInitPhase(observer, InitPhaseAutoProgressOff, func() error {
		shell.EndAutoProgress()
		return nil
	}); err != nil {
		return "", err
	}
	if options.computeVars {
		markInternalProgressStatPhaseComplete()
	}

	result := ""
	if err := runInitPhase(observer, InitPhaseOutputForm, func() error {
		result = shell.DotFile.String()
		if strings.Contains(strings.ToLower(result), "aliae error:") {
			return errInitFailed
		}
		return nil
	}); err != nil {
		return "", err
	}
	if options.computeVars {
		markInternalProgressOutputFormulated()
		markInternalProgressReadyToOutput()
	}

	return result, nil
}
