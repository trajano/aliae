package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Scripts []*Script

type ScriptType string

const (
	ShellScript  ScriptType = "shell"
	PythonScript ScriptType = "python"
	PerlScript   ScriptType = "perl"
)

type Script struct {
	Value  Template   `yaml:"value"`
	Type   ScriptType `yaml:"type"`
	If     If         `yaml:"if"`
	Weight float64    `yaml:"weight"`
}

func (s *Script) String() string {
	script := s.Value.Parse()
	if len(s.Type) == 0 || s.Type == ShellScript {
		return string(script)
	}

	switch s.Type {
	case ShellScript:
		return string(script)
	case PythonScript:
		return inlineInterpreterScript("python", "-c", string(script))
	case PerlScript:
		return inlineInterpreterScript("perl", "-e", string(script))
	default:
		return ""
	}
}

func (s *Script) effectiveWeight() float64 {
	if s.Weight <= 0 {
		return 1
	}

	return s.Weight
}

func (s Scripts) Render() {
	if len(s) == 0 {
		return
	}

	first := true
	for _, script := range s {
		if script.If.Ignore() {
			continue
		}

		scriptBlock := script.String()
		if len(scriptBlock) == 0 {
			advanceAutoProgress(script.effectiveWeight())
			continue
		}

		if first && dotFileHasRenderableContent() {
			if !dotFileEndsWithNewline() {
				DotFile.WriteString("\n")
			}
			DotFile.WriteString("\n")
		} else if !dotFileEndsWithNewline() {
			DotFile.WriteString("\n")
		}
		DotFile.WriteString(scriptBlock)

		first = false
		advanceAutoProgress(script.effectiveWeight())
	}
}

func inlineInterpreterScript(executable, switchName, script string) string {
	formatted, ok := formatString(Template(script)).(string)
	if !ok {
		return ""
	}

	command := fmt.Sprintf("%s %s %s", executable, switchName, formatted)

	if context.Current.Shell == CMD {
		luaFormatted, ok := formatString(command).(string)
		if !ok {
			return ""
		}
		return fmt.Sprintf("os.execute(%s)", luaFormatted)
	}

	return command
}
