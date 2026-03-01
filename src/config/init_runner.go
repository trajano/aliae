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

	emitInitVisits(aliae, observer)

	if err := runInitPhase(observer, InitPhaseAutoProgressOn, func() error {
		shell.StartAutoProgress(aliae.autoProgressConfig())
		return nil
	}); err != nil {
		return "", err
	}
	shell.AdvanceAutoProgress(aliae.progressVarWeight())

	if err := runInitRenderPhases(aliae, observer); err != nil {
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

func runInitRenderPhases(aliae *Aliae, observer InitObserver) error {
	renderPhases := []struct {
		run   func()
		phase InitPhase
	}{
		{run: func() { aliae.Envs.Render() }, phase: InitPhaseRenderEnv},
		{run: func() { aliae.Paths.Render() }, phase: InitPhaseRenderPath},
		{run: func() { aliae.CDPaths.Render() }, phase: InitPhaseRenderCDPath},
		{run: func() { aliae.Aliae.Render() }, phase: InitPhaseRenderAlias},
		{run: func() { aliae.Links.Render() }, phase: InitPhaseRenderLink},
		{run: func() { aliae.Scripts.Render() }, phase: InitPhaseRenderScript},
	}

	for _, phase := range renderPhases {
		step := phase
		if err := runInitPhase(observer, step.phase, func() error {
			step.run()
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func emitInitVisits(aliae *Aliae, observer InitObserver) {
	if observer == nil || aliae == nil {
		return
	}

	WalkConfig(aliae, ConfigVisitorFuncs{
		OnExtend: func(item *Extend) {
			key := strings.TrimSpace(item.Path)
			if key == "" {
				key = strings.TrimSpace(item.Dir)
			}
			runInitVisit(observer, InitSectionExtend, key)
		},
		OnVar: func(item *Var) {
			runInitVisit(observer, InitSectionVar, item.Name)
		},
		OnEnv: func(item *shell.Env) {
			runInitVisit(observer, InitSectionEnv, item.Name)
		},
		OnPath: func(item *shell.Path) {
			runInitVisit(observer, InitSectionPath, firstLine(string(item.Value)))
		},
		OnCDPath: func(item *shell.CDPath) {
			runInitVisit(observer, InitSectionCDPath, firstLine(string(item.Value)))
		},
		OnAlias: func(item *shell.Alias) {
			runInitVisit(observer, InitSectionAlias, item.Name)
		},
		OnLink: func(item *shell.Link) {
			runInitVisit(observer, InitSectionLink, firstLine(string(item.Name)))
		},
		OnScript: func(item *shell.Script) {
			runInitVisit(observer, InitSectionScript, firstLine(string(item.Value)))
		},
	})
}

func firstLine(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	if idx := strings.Index(trimmed, "\n"); idx >= 0 {
		return strings.TrimSpace(trimmed[:idx])
	}

	return trimmed
}
