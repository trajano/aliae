package config

import (
	"time"
)

type InitBenchmarkStep struct {
	Name     string
	Duration time.Duration
}

type InitVisitBenchmark struct {
	Section  InitSection
	Count    int
	Duration time.Duration
}

func BenchmarkInit(configPath, sh string) ([]InitBenchmarkStep, []InitVisitBenchmark, error) {
	observer := &initBenchmarkObserver{
		steps: make([]InitBenchmarkStep, 0, 12),
		visit: map[InitSection]InitVisitBenchmark{},
	}

	if _, err := runInitWithObserver(configPath, sh, observer, initRunOptions{
		computeVars: false,
		primeState:  false,
	}); err != nil {
		return nil, nil, err
	}

	return observer.steps, observer.visitSummaries(), nil
}

var errInitFailed = &initBenchmarkError{message: "init benchmark failed"}

type initBenchmarkError struct {
	message string
}

func (e *initBenchmarkError) Error() string {
	return e.message
}

type initBenchmarkObserver struct {
	visit map[InitSection]InitVisitBenchmark
	steps []InitBenchmarkStep
}

func (o *initBenchmarkObserver) OnInitPhaseStart(_ InitPhase) {}

func (o *initBenchmarkObserver) OnInitPhaseEnd(phase InitPhase, duration time.Duration, err error) {
	if err != nil {
		return
	}

	o.steps = append(o.steps, InitBenchmarkStep{
		Name:     string(phase),
		Duration: duration,
	})
}

func (o *initBenchmarkObserver) OnInitVisitStart(_ InitSection, _ string) {}

func (o *initBenchmarkObserver) OnInitVisitEnd(section InitSection, _ string, duration time.Duration) {
	current := o.visit[section]
	current.Section = section
	current.Count++
	current.Duration += duration
	o.visit[section] = current
}

func (o *initBenchmarkObserver) visitSummaries() []InitVisitBenchmark {
	order := []InitSection{
		InitSectionExtend,
		InitSectionVar,
		InitSectionEnv,
		InitSectionPath,
		InitSectionCDPath,
		InitSectionAlias,
		InitSectionLink,
		InitSectionScript,
	}

	summaries := make([]InitVisitBenchmark, 0, len(order))
	for _, section := range order {
		item, ok := o.visit[section]
		if !ok || item.Count == 0 {
			continue
		}
		summaries = append(summaries, item)
	}

	return summaries
}
