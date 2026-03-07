package config

import "github.com/jandedobbeleer/aliae/src/shell"

func AutoProgressConfig(a *Aliae) shell.AutoProgressConfig {
	return a.autoProgressConfig()
}

func (a *Aliae) autoProgressConfig() shell.AutoProgressConfig {
	if a == nil || !a.Progress.Enabled {
		return shell.AutoProgressConfig{}
	}

	start := a.Progress.StartPercentage + a.Progress.Internal
	if start > 100 {
		start = 100
	}

	return shell.AutoProgressConfig{
		Enabled:         true,
		StartPercentage: start,
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
	WalkConfig(a, ConfigVisitorFuncs{
		OnVar: func(variable *Var) {
			if variable.If.Ignore() {
				return
			}
			total += 1
		},
		OnAlias: func(alias *shell.Alias) {
			if alias.If.Ignore() {
				return
			}
			total += 1
		},
		OnEnv: func(env *shell.Env) {
			if env.If.Ignore() {
				return
			}
			total += 1
		},
		OnPath: func(path *shell.Path) {
			if path.If.Ignore() {
				return
			}
			total += 1
		},
		OnCDPath: func(cdpath *shell.CDPath) {
			if cdpath.If.Ignore() {
				return
			}
			total += 1
		},
		OnLink: func(link *shell.Link) {
			if link.If.Ignore() {
				return
			}
			total += 1
		},
		OnScript: func(script *shell.Script) {
			if script.If.Ignore() {
				return
			}
			total += scriptWeight(script)
		},
	})

	return total
}

func scriptWeight(script *shell.Script) float64 {
	if script.Weight <= 0 {
		return 1
	}

	return script.Weight
}

func ProgressVarWeight(a *Aliae) float64 {
	return a.progressVarWeight()
}

func (a *Aliae) progressVarWeight() float64 {
	if a == nil {
		return 0
	}

	total := 0.0
	for _, variable := range a.Vars {
		if variable == nil || variable.If.Ignore() {
			continue
		}
		total += 1
	}

	return total
}
