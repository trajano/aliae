package init

import "time"

type Phase string

const (
	PhaseContextInit     Phase = "init.context_init"
	PhaseLoadConfig      Phase = "init.load_config"
	PhaseEvaluateVars    Phase = "init.evaluate_vars"
	PhaseStatePrime      Phase = "init.state_prime"
	PhaseAutoProgressOn  Phase = "init.autoprogress_start"
	PhaseRenderEnv       Phase = "init.render_env"
	PhaseRenderPath      Phase = "init.render_path"
	PhaseRenderCDPath    Phase = "init.render_cdpath"
	PhaseRenderAlias     Phase = "init.render_alias"
	PhaseRenderLink      Phase = "init.render_link"
	PhaseRenderScript    Phase = "init.render_script"
	PhaseAutoProgressOff Phase = "init.autoprogress_end"
	PhaseOutputForm      Phase = "init.output_form"
)

type Observer interface {
	OnPhaseStart(phase Phase)
	OnPhaseEnd(phase Phase, duration time.Duration, err error)
	OnVisitStart(section Section, key string)
	OnVisitEnd(section Section, key string, duration time.Duration)
}

type Section string

const (
	SectionExtend     Section = "extends"
	SectionVar        Section = "var"
	SectionEnv        Section = "env"
	SectionPath       Section = "path"
	SectionCDPath     Section = "cdpath"
	SectionAlias      Section = "alias"
	SectionLink       Section = "link"
	SectionScript     Section = "script"
	SectionStateCheck Section = "state_check"
)

func OrderedSections() []Section {
	return []Section{
		SectionExtend,
		SectionVar,
		SectionEnv,
		SectionPath,
		SectionCDPath,
		SectionAlias,
		SectionLink,
		SectionScript,
		SectionStateCheck,
	}
}

func runPhase(observer Observer, phase Phase, run func() error) error {
	if observer != nil {
		observer.OnPhaseStart(phase)
	}

	start := time.Now()
	err := run()

	if observer != nil {
		observer.OnPhaseEnd(phase, time.Since(start), err)
	}

	return err
}

func runVisit(observer Observer, section Section, key string) {
	if observer == nil {
		return
	}

	observer.OnVisitStart(section, key)
	start := time.Now()
	observer.OnVisitEnd(section, key, time.Since(start))
}
