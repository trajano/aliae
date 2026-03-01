package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jandedobbeleer/aliae/src/shell"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

type validationCollector struct {
	lineResolver *yamlLineResolver
	errors       []string
}

func newValidationCollector(lineResolver *yamlLineResolver) *validationCollector {
	return &validationCollector{
		lineResolver: lineResolver,
		errors:       make([]string, 0),
	}
}

func (c *validationCollector) add(path, message string) {
	c.errors = append(c.errors, c.lineResolver.annotate(path, message))
}

func (c *validationCollector) err(prefix string) error {
	if len(c.errors) == 0 {
		return nil
	}

	slices.Sort(c.errors)
	return fmt.Errorf("%s:\n- %s", prefix, strings.Join(c.errors, "\n- "))
}

type ifExpressionValidationVisitor struct{}

func (v ifExpressionValidationVisitor) Visit(aliae *Aliae, collector *validationCollector) {
	if aliae == nil {
		return
	}

	index := buildValidationIndex(aliae)
	WalkConfig(aliae, ConfigVisitorFuncs{
		OnVar: func(variable *Var) {
			if err := variable.If.Validate(); err != nil {
				i := index.varIndex[variable]
				collector.add(fmt.Sprintf("var.%d.if", i), fmt.Sprintf("var[%d].if: %s", i, err))
			}
		},
		OnEnv: func(env *shell.Env) {
			if err := env.If.Validate(); err != nil {
				i := index.envIndex[env]
				collector.add(fmt.Sprintf("env.%d.if", i), fmt.Sprintf("env[%d].if: %s", i, err))
			}
		},
		OnPath: func(path *shell.Path) {
			if err := path.If.Validate(); err != nil {
				i := index.pathIndex[path]
				collector.add(fmt.Sprintf("path.%d.if", i), fmt.Sprintf("path[%d].if: %s", i, err))
			}
		},
		OnCDPath: func(path *shell.CDPath) {
			if err := path.If.Validate(); err != nil {
				i := index.cdpathIndex[path]
				collector.add(fmt.Sprintf("cdpath.%d.if", i), fmt.Sprintf("cdpath[%d].if: %s", i, err))
			}
		},
		OnAlias: func(alias *shell.Alias) {
			if err := alias.If.Validate(); err != nil {
				i := index.aliasIndex[alias]
				collector.add(fmt.Sprintf("alias.%d.if", i), fmt.Sprintf("alias[%d].if: %s", i, err))
			}
		},
		OnLink: func(link *shell.Link) {
			if err := link.If.Validate(); err != nil {
				i := index.linkIndex[link]
				collector.add(fmt.Sprintf("link.%d.if", i), fmt.Sprintf("link[%d].if: %s", i, err))
			}
		},
		OnScript: func(script *shell.Script) {
			if err := script.If.Validate(); err != nil {
				i := index.scriptIndex[script]
				collector.add(fmt.Sprintf("script.%d.if", i), fmt.Sprintf("script[%d].if: %s", i, err))
			}
		},
	})
}

type scriptWeightValidationVisitor struct{}

func (v scriptWeightValidationVisitor) Visit(aliae *Aliae, collector *validationCollector) {
	if aliae == nil {
		return
	}

	index := buildValidationIndex(aliae)
	WalkConfig(aliae, ConfigVisitorFuncs{
		OnScript: func(script *shell.Script) {
			if script.Weight == 0 {
				return
			}

			if script.Weight < 0 {
				i := index.scriptIndex[script]
				collector.add(
					fmt.Sprintf("script.%d.weight", i),
					fmt.Sprintf("script[%d].weight must be greater than 0", i),
				)
			}
		},
	})
}

type scriptStateValidationVisitor struct{}

func (v scriptStateValidationVisitor) Visit(aliae *Aliae, collector *validationCollector) {
	if aliae == nil {
		return
	}

	index := buildValidationIndex(aliae)
	seenStateFiles := make(map[string]int)

	WalkConfig(aliae, ConfigVisitorFuncs{
		OnScript: func(script *shell.Script) {
			if len(script.State.File) == 0 {
				return
			}

			i := index.scriptIndex[script]
			stateFile := strings.TrimSpace(string(script.State.File))
			filePath := fmt.Sprintf("script.%d.state.file", i)
			if !aliaeState.IsValidFileName(stateFile) {
				collector.add(filePath, fmt.Sprintf("script[%d].state.file must be a file name only (no path separators)", i))
			}

			if previousIndex, exists := seenStateFiles[stateFile]; exists {
				collector.add(filePath, fmt.Sprintf("script[%d].state.file duplicates script[%d].state.file", i, previousIndex))
			} else {
				seenStateFiles[stateFile] = i
			}

			if len(script.State.RunEvery) > 0 {
				runEveryPath := fmt.Sprintf("script.%d.state.runEvery", i)
				runEvery, err := time.ParseDuration(strings.TrimSpace(script.State.RunEvery))
				if err != nil {
					collector.add(runEveryPath, fmt.Sprintf("script[%d].state.runEvery must be a valid duration", i))
				} else if runEvery <= 0 {
					collector.add(runEveryPath, fmt.Sprintf("script[%d].state.runEvery must be greater than 0", i))
				}
			}
		},
	})
}

type progressValidationVisitor struct{}

func (v progressValidationVisitor) Visit(aliae *Aliae, collector *validationCollector) {
	if aliae == nil || !aliae.Progress.Enabled {
		return
	}

	start := aliae.Progress.StartPercentage
	internal := aliae.Progress.Internal
	effectiveStart := start + internal
	end := aliae.Progress.EndPercentage.Value

	if start < 0 || start > 100 {
		collector.add("progress.start_percentage", "progress.start_percentage must be between 0 and 100")
	}

	if internal < 0 || internal > 100 {
		collector.add("progress.internal", "progress.internal must be between 0 and 100")
	}

	if effectiveStart > 100 {
		collector.add("progress.internal", "progress.start_percentage + progress.internal must be less than or equal to 100")
	}

	if end < 0 || end > 100 {
		collector.add("progress.end_percentage", "progress.end_percentage must be between 0 and 100")
	}

	if !aliae.Progress.EndPercentage.Reset && end <= effectiveStart {
		collector.add("progress.end_percentage", "progress.end_percentage must be greater than progress.start_percentage + progress.internal")
	}
}

type extendsIfValidationVisitor struct{}

func (v extendsIfValidationVisitor) Visit(aliae *Aliae, collector *validationCollector) {
	if aliae == nil {
		return
	}

	for i := range aliae.Extends {
		item := aliae.Extends[i]
		if err := item.If.Validate(); err != nil {
			collector.add(fmt.Sprintf("extends.%d.if", i), fmt.Sprintf("extends[%d].if: %s", i, err))
		}
	}
}

type validationIndex struct {
	aliasIndex  map[*shell.Alias]int
	varIndex    map[*Var]int
	envIndex    map[*shell.Env]int
	pathIndex   map[*shell.Path]int
	cdpathIndex map[*shell.CDPath]int
	scriptIndex map[*shell.Script]int
	linkIndex   map[*shell.Link]int
}

func buildValidationIndex(aliae *Aliae) validationIndex {
	index := validationIndex{
		aliasIndex:  map[*shell.Alias]int{},
		varIndex:    map[*Var]int{},
		envIndex:    map[*shell.Env]int{},
		pathIndex:   map[*shell.Path]int{},
		cdpathIndex: map[*shell.CDPath]int{},
		scriptIndex: map[*shell.Script]int{},
		linkIndex:   map[*shell.Link]int{},
	}

	for i, item := range aliae.Aliae {
		if item != nil {
			index.aliasIndex[item] = i
		}
	}
	for i, item := range aliae.Vars {
		if item != nil {
			index.varIndex[item] = i
		}
	}
	for i, item := range aliae.Envs {
		if item != nil {
			index.envIndex[item] = i
		}
	}
	for i, item := range aliae.Paths {
		if item != nil {
			index.pathIndex[item] = i
		}
	}
	for i, item := range aliae.CDPaths {
		if item != nil {
			index.cdpathIndex[item] = i
		}
	}
	for i, item := range aliae.Scripts {
		if item != nil {
			index.scriptIndex[item] = i
		}
	}
	for i, item := range aliae.Links {
		if item != nil {
			index.linkIndex[item] = i
		}
	}

	return index
}
