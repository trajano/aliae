package shell

import "math"

type AutoProgressConfig struct {
	StartPercentage float64
	EndPercentage   float64
	TotalWeight     float64
	Enabled         bool
	ResetAtEnd      bool
}

type autoProgressTracker struct {
	startPercentage float64
	endPercentage   float64
	totalWeight     float64
	processedWeight float64
	lastPercentage  int
	resetAtEnd      bool
	lastSet         bool
}

var autoProgress *autoProgressTracker

func StartAutoProgress(cfg AutoProgressConfig) {
	if !cfg.Enabled || cfg.TotalWeight <= 0 {
		autoProgress = nil
		return
	}

	autoProgress = &autoProgressTracker{
		startPercentage: cfg.StartPercentage,
		endPercentage:   cfg.EndPercentage,
		totalWeight:     cfg.TotalWeight,
		resetAtEnd:      cfg.ResetAtEnd,
	}

	writeAutoProgress(int(math.Floor(cfg.StartPercentage)))
}

func advanceAutoProgress(weight float64) {
	if autoProgress == nil {
		return
	}

	if weight <= 0 {
		weight = 1
	}

	autoProgress.processedWeight += weight
	if autoProgress.processedWeight > autoProgress.totalWeight {
		autoProgress.processedWeight = autoProgress.totalWeight
	}

	percentage := autoProgress.startPercentage
	span := autoProgress.endPercentage - autoProgress.startPercentage
	if span != 0 {
		percentage += (autoProgress.processedWeight / autoProgress.totalWeight) * span
	}

	writeAutoProgress(int(math.Floor(percentage)))
}

func EndAutoProgress() {
	if autoProgress == nil {
		return
	}

	if autoProgress.resetAtEnd {
		writeAutoProgress(99)
		writeDotFileProgress("reset")
		autoProgress = nil
		return
	}

	writeAutoProgress(int(math.Floor(autoProgress.endPercentage)))
	autoProgress = nil
}

func writeAutoProgress(percentage int) {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	if autoProgress.resetAtEnd && percentage >= 100 {
		percentage = 99
	}

	if autoProgress.lastSet && autoProgress.lastPercentage == percentage {
		return
	}

	autoProgress.lastSet = true
	autoProgress.lastPercentage = percentage
	writeDotFileProgress(percentage)
}

func writeDotFileProgress(value any) {
	command := progress(value)
	if len(command) == 0 {
		return
	}

	if DotFile.Len() > 0 && !dotFileEndsWithNewline() {
		DotFile.WriteString("\n")
	}

	DotFile.WriteString(command)
	DotFile.WriteString("\n")
}

func dotFileEndsWithNewline() bool {
	text := DotFile.String()
	return len(text) > 0 && text[len(text)-1] == '\n'
}

func dotFileHasRenderableContent() bool {
	text := DotFile.String()
	if len(text) == 0 {
		return false
	}

	lines := splitAndTrimLines(text)
	for _, line := range lines {
		if len(line) == 0 || isProgressLine(line) {
			continue
		}

		return true
	}

	return false
}

func splitAndTrimLines(text string) []string {
	lines := []string{}
	current := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, trimSpace(current))
			current = ""
			continue
		}

		if r == '\r' {
			continue
		}

		current += string(r)
	}

	lines = append(lines, trimSpace(current))
	return lines
}

func trimSpace(value string) string {
	start := 0
	end := len(value)

	for start < end && (value[start] == ' ' || value[start] == '\t') {
		start++
	}

	for end > start && (value[end-1] == ' ' || value[end-1] == '\t') {
		end--
	}

	return value[start:end]
}

func isProgressLine(line string) bool {
	return len(line) > 0 && (hasPrefix(line, `printf '\033]9;4;`) || hasPrefix(line, `[Console]::Out.Write("$([char]27)]9;4;`))
}

func hasPrefix(text, prefix string) bool {
	if len(prefix) > len(text) {
		return false
	}

	return text[:len(prefix)] == prefix
}
