package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type Vars []*Var

type Var struct {
	Name        string         `yaml:"name"`
	Value       shell.Template `yaml:"value"`
	If          shell.If       `yaml:"if"`
	ifEvaluated bool           `yaml:"-"`
	ifIgnored   bool           `yaml:"-"`
}

func (a *Aliae) ComputeVars() error {
	return a.computeVars(nil)
}

func (a *Aliae) computeVars(lineResolver *yamlLineResolver) error {
	if a == nil {
		return nil
	}

	ctx := ensureTemplateRuntime()

	if ctx.Var == nil {
		ctx.Var = map[string]any{}
	}

	if err := validateVarDefinitions(a.Vars, lineResolver); err != nil {
		return err
	}

	computed := make(map[string]any, len(a.Vars))
	ctx.Var = computed

	for _, variable := range a.Vars {
		if variable == nil || variable.Ignore() {
			continue
		}

		computed[variable.Name] = evaluateVarValue(variable.Value)
	}

	ctx.Var = computed
	MarkInitInternalProgressVarsComputed()
	return nil
}

func (v *Var) Ignore() bool {
	if v.ifEvaluated {
		return v.ifIgnored
	}

	v.ifIgnored = v.If.Ignore()
	v.ifEvaluated = true
	return v.ifIgnored
}

func evaluateVarValue(value shell.Template) any {
	text := strings.TrimSpace(string(value))
	if len(text) == 0 {
		return ""
	}

	if strings.Contains(text, "{{") && strings.Contains(text, "}}") {
		return normalizeVarValue(value.String())
	}

	// Allow `var.value` to behave like `if`, where bare expressions are valid.
	parsed := shell.Template(fmt.Sprintf("{{ %s }}", text)).String()
	if parsed == fmt.Sprintf("{{ %s }}", text) {
		return text
	}

	return normalizeVarValue(parsed)
}

func normalizeVarValue(value string) any {
	trimmed := strings.TrimSpace(value)
	parsedBool, err := strconv.ParseBool(trimmed)
	if err == nil {
		return parsedBool
	}

	return value
}

func ensureTemplateRuntime() *context.Runtime {
	if current := context.GetCurrent(); current != nil {
		shell.SetRuntime(current)
		return current
	}

	current := &context.Runtime{
		Shell: shell.BASH,
		OS:    context.LINUX,
		Home:  context.Home(),
		Env:   map[string]string{},
		Var:   map[string]any{},
	}
	context.SetCurrent(current)
	shell.SetRuntime(current)

	return current
}

func validateVarDefinitions(vars Vars, lineResolver *yamlLineResolver) error {
	validationErrors := make([]string, 0)

	for i, variable := range vars {
		if variable == nil {
			continue
		}

		if strings.Contains(string(variable.Value), ".Var") {
			path := fmt.Sprintf("var.%d.value", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("var[%d].value cannot reference .Var", i)))
		}

		if strings.Contains(string(variable.If), ".Var") {
			path := fmt.Sprintf("var.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("var[%d].if cannot reference .Var", i)))
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	slices.Sort(validationErrors)
	return fmt.Errorf("config var validation failed:\n- %s", strings.Join(validationErrors, "\n- "))
}
