package config

import (
	"time"
)

type InitBenchmarkStep struct {
	Name     string
	Duration time.Duration
}

func BenchmarkInit(configPath, sh string) ([]InitBenchmarkStep, error) {
	observer := &initBenchmarkObserver{
		steps: make([]InitBenchmarkStep, 0, 12),
	}

	if _, err := runInitWithObserver(configPath, sh, observer, initRunOptions{
		computeVars: false,
		primeState:  false,
	}); err != nil {
		return nil, err
	}

	return observer.steps, nil
}

var errInitFailed = &initBenchmarkError{message: "init benchmark failed"}

type initBenchmarkError struct {
	message string
}

func (e *initBenchmarkError) Error() string {
	return e.message
}

type initBenchmarkObserver struct {
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

func (o *initBenchmarkObserver) OnInitVisitEnd(_ InitSection, _ string, _ time.Duration) {}
