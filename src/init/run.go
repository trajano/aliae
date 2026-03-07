package init

import (
	"strings"
	"time"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type runOptions struct {
	computeVars bool
	primeState  bool
}

func runWithObserver(configPath, sh string, observer Observer, options runOptions) (string, error) {
	shell.DotFile.Reset()
	defer shell.DotFile.Reset()

	var aliae *cfg.Aliae
	if err := runPhase(observer, PhaseContextInit, func() error {
		context.Init(sh)
		return nil
	}); err != nil {
		return "", err
	}

	if err := runPhase(observer, PhaseLoadConfig, func() error {
		var err error
		if options.computeVars {
			aliae, err = cfg.LoadConfig(configPath)
		} else {
			aliae, err = cfg.LoadConfigWithoutVars(configPath)
		}
		return err
	}); err != nil {
		return "", err
	}

	if options.computeVars {
		cfg.MarkInitInternalProgressConfigValidated()
	} else {
		if err := runPhase(observer, PhaseEvaluateVars, func() error {
			return aliae.ComputeVars()
		}); err != nil {
			return "", err
		}
	}

	if options.primeState {
		stateChecks := aliae.Scripts.PrimeState(time.Now())
		cfg.MarkInitInternalProgressStateChecksComplete(stateChecks)
	} else {
		if err := runPhase(observer, PhaseStatePrime, func() error {
			_ = aliae.Scripts.PrimeState(time.Now())
			return nil
		}); err != nil {
			return "", err
		}
	}

	emitVisits(aliae, observer)

	if err := runPhase(observer, PhaseAutoProgressOn, func() error {
		shell.StartAutoProgress(cfg.AutoProgressConfig(aliae))
		return nil
	}); err != nil {
		return "", err
	}
	shell.AdvanceAutoProgress(cfg.ProgressVarWeight(aliae))

	if err := runRenderPhases(aliae, observer); err != nil {
		return "", err
	}

	if err := runPhase(observer, PhaseAutoProgressOff, func() error {
		shell.EndAutoProgress()
		return nil
	}); err != nil {
		return "", err
	}
	if options.computeVars {
		cfg.MarkInitInternalProgressStatPhaseComplete()
	}

	result := ""
	if err := runPhase(observer, PhaseOutputForm, func() error {
		result = shell.DotFile.String()
		if strings.Contains(strings.ToLower(result), "aliae error:") {
			return errInitFailed
		}
		return nil
	}); err != nil {
		return "", err
	}
	if options.computeVars {
		cfg.MarkInitInternalProgressOutputFormulated()
		cfg.MarkInitInternalProgressReadyToOutput()
	}

	return result, nil
}

func runRenderPhases(aliae *cfg.Aliae, observer Observer) error {
	if aliae == nil {
		return nil
	}

	renderPhases := []struct {
		renderer interface{ Render() }
		phase    Phase
	}{
		{renderer: aliae.Envs, phase: PhaseRenderEnv},
		{renderer: aliae.Paths, phase: PhaseRenderPath},
		{renderer: aliae.CDPaths, phase: PhaseRenderCDPath},
		{renderer: aliae.Aliae, phase: PhaseRenderAlias},
		{renderer: aliae.Links, phase: PhaseRenderLink},
		{renderer: aliae.Scripts, phase: PhaseRenderScript},
	}

	if observer == nil {
		for _, step := range renderPhases {
			step.renderer.Render()
		}
		return nil
	}

	for _, step := range renderPhases {
		current := step
		if err := runPhase(observer, current.phase, func() error {
			current.renderer.Render()
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func emitVisits(aliae *cfg.Aliae, observer Observer) {
	if observer == nil || aliae == nil {
		return
	}

	cfg.WalkConfig(aliae, cfg.ConfigVisitorFuncs{
		OnExtend: func(item *cfg.Extend) {
			key := strings.TrimSpace(item.Path)
			if key == "" {
				key = strings.TrimSpace(item.Dir)
			}
			runVisit(observer, SectionExtend, key)
		},
		OnVar: func(item *cfg.Var) {
			runVisit(observer, SectionVar, item.Name)
		},
		OnEnv: func(item *shell.Env) {
			runVisit(observer, SectionEnv, item.Name)
		},
		OnPath: func(item *shell.Path) {
			runVisit(observer, SectionPath, firstLine(string(item.Value)))
		},
		OnCDPath: func(item *shell.CDPath) {
			runVisit(observer, SectionCDPath, firstLine(string(item.Value)))
		},
		OnAlias: func(item *shell.Alias) {
			runVisit(observer, SectionAlias, item.Name)
		},
		OnLink: func(item *shell.Link) {
			runVisit(observer, SectionLink, firstLine(string(item.Name)))
		},
		OnScript: func(item *shell.Script) {
			runVisit(observer, SectionScript, firstLine(string(item.Value)))
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
