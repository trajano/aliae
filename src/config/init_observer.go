package config

import "time"

type InitPhase string

const (
	InitPhaseContextInit     InitPhase = "init.context_init"
	InitPhaseLoadConfig      InitPhase = "init.load_config"
	InitPhaseEvaluateVars    InitPhase = "init.evaluate_vars"
	InitPhaseStatePrime      InitPhase = "init.state_prime"
	InitPhaseAutoProgressOn  InitPhase = "init.autoprogress_start"
	InitPhaseRenderEnv       InitPhase = "init.render_env"
	InitPhaseRenderPath      InitPhase = "init.render_path"
	InitPhaseRenderCDPath    InitPhase = "init.render_cdpath"
	InitPhaseRenderAlias     InitPhase = "init.render_alias"
	InitPhaseRenderLink      InitPhase = "init.render_link"
	InitPhaseRenderScript    InitPhase = "init.render_script"
	InitPhaseAutoProgressOff InitPhase = "init.autoprogress_end"
	InitPhaseOutputForm      InitPhase = "init.output_form"
)

type InitObserver interface {
	OnInitPhaseStart(phase InitPhase)
	OnInitPhaseEnd(phase InitPhase, duration time.Duration, err error)
}

func runInitPhase(observer InitObserver, phase InitPhase, run func() error) error {
	if observer != nil {
		observer.OnInitPhaseStart(phase)
	}

	start := time.Now()
	err := run()

	if observer != nil {
		observer.OnInitPhaseEnd(phase, time.Since(start), err)
	}

	return err
}
