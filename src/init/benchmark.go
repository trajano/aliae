package init

import (
	"time"

	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

type BenchmarkStep struct {
	Name     string
	Duration time.Duration
}

type VisitBenchmark struct {
	Section  Section
	Count    int
	Duration time.Duration
}

func Benchmark(configPath, sh string) ([]BenchmarkStep, []VisitBenchmark, error) {
	observer := &benchmarkObserver{
		steps: make([]BenchmarkStep, 0, 12),
		visit: map[Section]VisitBenchmark{},
	}
	aliaeState.ResetMetrics()
	aliaeState.EnableMetrics(true)
	defer aliaeState.EnableMetrics(false)

	if _, err := runWithObserver(configPath, sh, observer, runOptions{
		computeVars: false,
		primeState:  false,
	}); err != nil {
		return nil, nil, err
	}

	metrics := aliaeState.SnapshotMetrics()
	if metrics.ShouldRunCount > 0 {
		observer.visit[SectionStateCheck] = VisitBenchmark{
			Section:  SectionStateCheck,
			Count:    int(metrics.ShouldRunCount),
			Duration: metrics.ShouldRunTime,
		}
	}

	return observer.steps, observer.visitSummaries(), nil
}

var errInitFailed = &benchmarkError{message: "init benchmark failed"}

type benchmarkError struct {
	message string
}

func (e *benchmarkError) Error() string {
	return e.message
}

type benchmarkObserver struct {
	visit map[Section]VisitBenchmark
	steps []BenchmarkStep
}

func (o *benchmarkObserver) OnPhaseStart(_ Phase) {}

func (o *benchmarkObserver) OnPhaseEnd(phase Phase, duration time.Duration, err error) {
	if err != nil {
		return
	}

	o.steps = append(o.steps, BenchmarkStep{
		Name:     string(phase),
		Duration: duration,
	})
}

func (o *benchmarkObserver) OnVisitStart(_ Section, _ string) {}

func (o *benchmarkObserver) OnVisitEnd(section Section, _ string, duration time.Duration) {
	current := o.visit[section]
	current.Section = section
	current.Count++
	current.Duration += duration
	o.visit[section] = current
}

func (o *benchmarkObserver) visitSummaries() []VisitBenchmark {
	sections := OrderedSections()
	summaries := make([]VisitBenchmark, 0, len(sections))
	for _, section := range sections {
		item, ok := o.visit[section]
		if !ok || item.Count == 0 {
			continue
		}
		summaries = append(summaries, item)
	}

	return summaries
}
