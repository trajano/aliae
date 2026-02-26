package config

import "github.com/jandedobbeleer/aliae/src/shell"

func (a *Aliae) autoProgressConfig() shell.AutoProgressConfig {
	if a == nil || !a.Progress.Enabled {
		return shell.AutoProgressConfig{}
	}

	return shell.AutoProgressConfig{
		Enabled:         true,
		StartPercentage: a.Progress.StartPercentage,
		EndPercentage:   a.Progress.EndPercentage.Value,
		ResetAtEnd:      a.Progress.EndPercentage.Reset,
		TotalWeight:     a.progressTotalWeight(),
	}
}

func (a *Aliae) progressTotalWeight() float64 {
	if a == nil {
		return 0
	}

	total := 0.0

	for _, alias := range a.Aliae {
		if alias == nil || alias.If.Ignore() {
			continue
		}

		total += 1
	}

	for _, env := range a.Envs {
		if env == nil || env.If.Ignore() {
			continue
		}

		total += 1
	}

	for _, path := range a.Paths {
		if path == nil || path.If.Ignore() {
			continue
		}

		total += 1
	}

	for _, script := range a.Scripts {
		if script == nil || script.If.Ignore() {
			continue
		}

		total += scriptWeight(script)
	}

	return total
}

func scriptWeight(script *shell.Script) float64 {
	if script.Weight == nil {
		return 1
	}

	return *script.Weight
}
