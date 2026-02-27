package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type Vars []*Var

type Var struct {
	Name  string         `yaml:"name"`
	Value shell.Template `yaml:"value"`
	If    shell.If       `yaml:"if"`
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
		ctx.Var = map[string]string{}
	}

	if err := validateVarDefinitions(a.Vars, lineResolver); err != nil {
		return err
	}

	computed := make(map[string]string, len(a.Vars))
	ctx.Var = computed

	for _, variable := range a.Vars {
		if variable == nil || variable.If.Ignore() {
			continue
		}

		computed[variable.Name] = variable.Value.String()
	}

	ctx.Var = computed
	markInternalProgressVarsComputed()
	return nil
}

func ensureTemplateRuntime() *context.Runtime {
	if context.Current != nil {
		return context.Current
	}

	context.Current = &context.Runtime{
		Shell: shell.BASH,
		OS:    context.LINUX,
		Home:  context.Home(),
		Env:   map[string]string{},
		Var:   map[string]string{},
	}

	return context.Current
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
