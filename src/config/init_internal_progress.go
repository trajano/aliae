package config

import (
	"fmt"
	"io"
	"math"
	"os"
)

var (
	initProgressWriter io.Writer = os.Stderr
	internalProgress   *initInternalProgressReporter
)

type initInternalProgressReporter struct {
	enabled bool
	start   int
	end     int
	current int
	last    int
}

func SetInitProgressWriter(writer io.Writer) {
	if writer == nil {
		initProgressWriter = os.Stderr
		return
	}

	initProgressWriter = writer
}

func BeginInitInternalProgress(configPath string) {
	progress, err := loadRootProgress(resolveConfigPath(configPath))
	if err != nil || !progress.Enabled || progress.Internal <= 0 {
		internalProgress = nil
		return
	}

	start := clampPercentage(int(math.Floor(progress.StartPercentage)))
	end := clampPercentage(int(math.Floor(progress.StartPercentage + progress.Internal)))
	if end < start {
		end = start
	}

	reporter := &initInternalProgressReporter{
		enabled: true,
		start:   start,
		end:     end,
		current: start,
		last:    -1,
	}
	reporter.emit(start)
	internalProgress = reporter
}

func EndInitInternalProgress() {
	internalProgress = nil
}

func MarkInitInternalProgressDiscoveryComplete() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(1)
}

func MarkInitInternalProgressLinkedConfigLoaded() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(1)
}

func MarkInitInternalProgressConfigValidated() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(1)
}

func MarkInitInternalProgressVarsComputed() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(1)
}

func MarkInitInternalProgressStatPhaseComplete() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(2)
}

func MarkInitInternalProgressStateChecksComplete(count int) {
	if internalProgress == nil || count <= 0 {
		return
	}

	internalProgress.advance(count)
}

func MarkInitInternalProgressOutputFormulated() {
	if internalProgress == nil {
		return
	}

	internalProgress.advance(1)
}

func MarkInitInternalProgressReadyToOutput() {
	if internalProgress == nil {
		return
	}

	internalProgress.emit(internalProgress.end)
}

func (r *initInternalProgressReporter) advance(delta int) {
	if !r.enabled || delta <= 0 {
		return
	}

	next := r.current + delta
	if next > r.end {
		next = r.end
	}

	r.current = next
	r.emit(next)
}

func (r *initInternalProgressReporter) emit(percentage int) {
	if !r.enabled {
		return
	}

	value := clampPercentage(percentage)
	if r.last == value {
		return
	}

	r.last = value
	if initProgressWriter == nil {
		return
	}

	_, _ = fmt.Fprintf(initProgressWriter, "\x1b]9;4;1;%d\x07", value)
}

func clampPercentage(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}
