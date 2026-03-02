package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/registry"
)

type Envs []*Env

type Env struct {
	Value     any      `yaml:"value"`
	Name      string   `yaml:"name"`
	Delimiter Template `yaml:"delimiter"`
	If        If       `yaml:"if"`
	Type      EnvType  `yaml:"type"`
	template  string
	IsPath    bool `yaml:"isPath"`
	IfExists  bool `yaml:"ifExists"`
	Persist   bool `yaml:"persist"`
	parsed    bool
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ";")
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (e *Env) string() string {
	return newShellFactory().strategy().renderEnv(e)
}

func (e *Env) join() {
	if len(e.Delimiter) == 0 {
		return
	}

	text, OK := e.Value.(string)
	if !OK {
		return
	}

	splitted := strings.Split(text, "\n")
	splitted = filterEmpty(splitted)
	if len(splitted) == 1 {
		e.Value = splitted[0]
		return
	}

	for index, value := range splitted {
		splitted[index] = strings.TrimSpace(value)
	}

	delimiter := e.Delimiter.String()

	e.Value = strings.Join(splitted, delimiter)
}

func (e *Env) parse() {
	if e.parsed {
		return
	}

	e.parsed = true

	template := Template(toString(e.Value))
	e.Value = template.Parse().String()
	e.normalizePath()
	e.join()
}

func (e *Env) normalizePath() {
	if !e.IsPath || context.Current == nil || context.Current.OS != context.WINDOWS {
		return
	}

	value, ok := e.Value.(string)
	if !ok {
		return
	}

	lines := strings.Split(value, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}

		if context.Current.Shell == BASH && isMSYS2Environment() {
			if msysPath, isWindowsPath := windowsToMSYSPath(trimmed); isWindowsPath {
				trimmed = msysPath
			}

			lines[i] = strings.ReplaceAll(trimmed, `\`, "/")
			continue
		}

		if windowsPath, isMSYSPath := msysToWindowsPath(trimmed); isMSYSPath {
			trimmed = windowsPath
		}

		lines[i] = strings.ReplaceAll(trimmed, "/", `\`)
	}

	e.Value = strings.Join(lines, "\n")
}

func (e *Env) render() string {
	e.parse()

	script, err := parse(e.template, e)
	if err != nil {
		return err.Error()
	}

	return script
}

func (e Envs) Render() {
	e = e.filter()

	if len(e) == 0 {
		return
	}

	if dotFileHasRenderableContent() {
		DotFile.WriteString("\n\n")
	}

	if context.Current.Shell == NU {
		DotFile.WriteString(NuEnvBlockStart)
	}

	first := true
	for _, variable := range e {
		if !first {
			if !dotFileEndsWithNewline() {
				DotFile.WriteString("\n")
			}
		}

		DotFile.WriteString(variable.string())
		advanceAutoProgress(1)

		os.Setenv(variable.Name, toString(variable.Value))

		first = false
	}

	if context.Current.Shell == NU {
		DotFile.WriteString(NuEnvBlockEnd)
	}
}

func (e Envs) filter() Envs {
	var env Envs

	for _, variable := range e {
		if variable.If.Ignore() {
			continue
		}

		if variable.IfExists || variable.Persist || variable.IsPath {
			variable.parse()
		}

		if !variable.shouldExportPathValue() {
			continue
		}

		if variable.Persist {
			registry.PersistEnvironmentVariable(variable.Name, variable.Value)
		}

		env = append(env, variable)
	}

	return env
}

func (e *Env) shouldExportPathValue() bool {
	if !e.IsPath {
		return !e.IfExists
	}

	value, ok := e.Value.(string)
	if !ok {
		return false
	}

	lines := strings.Split(value, "\n")
	candidates := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		candidates = append(candidates, trimmed)
	}

	if len(candidates) == 0 {
		return false
	}

	autoIfExists := len(candidates) > 1
	if !e.IfExists && !autoIfExists {
		return true
	}

	for _, candidate := range candidates {
		if dirExists(candidate) {
			e.Value = candidate
			return true
		}
	}

	return false
}

func filterEmpty[S ~[]E, E string](s S) S {
	var cleaned S
	for _, a := range s {
		if len(a) == 0 {
			continue
		}
		cleaned = append(cleaned, a)
	}
	return cleaned
}

type EnvType string

const (
	String EnvType = "string"
	Array  EnvType = "array"
)
