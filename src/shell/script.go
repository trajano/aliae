package shell

import (
	"fmt"
	"strings"
	"time"

	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

type Scripts []*Script

type ScriptType string

const (
	ShellScript  ScriptType = "shell"
	PythonScript ScriptType = "python"
	PerlScript   ScriptType = "perl"
)

type Script struct {
	Value       Template              `yaml:"value"`
	Type        ScriptType            `yaml:"type"`
	If          If                    `yaml:"if"`
	State       ScriptState           `yaml:"state"`
	statePath   string                `yaml:"-"`
	stateFormat aliaeState.FileFormat `yaml:"-"`
	Weight      float64               `yaml:"weight"`

	statePrepared  bool `yaml:"-"`
	stateChecked   bool `yaml:"-"`
	stateShouldRun bool `yaml:"-"`
	ifFrozen       bool `yaml:"-"`
	ifIgnoreFrozen bool `yaml:"-"`
}

type ScriptState struct {
	File     Template              `yaml:"file"`
	RunEvery string                `yaml:"runEvery"`
	Format   aliaeState.FileFormat `yaml:"format"`
}

type ScriptStateReference struct {
	File     string
	Format   aliaeState.FileFormat
	RunEvery time.Duration
}

func (s *Script) stateReference() (*ScriptStateReference, error) {
	file := strings.TrimSpace(s.State.File.Parse().String())
	if len(file) == 0 {
		return nil, nil
	}

	runEvery := time.Duration(0)
	runEveryRaw := strings.TrimSpace(s.State.RunEvery)
	if len(runEveryRaw) > 0 {
		parsed, err := time.ParseDuration(runEveryRaw)
		if err != nil {
			return nil, err
		}
		runEvery = parsed
	}

	format := s.State.Format
	if len(format) == 0 {
		format = aliaeState.FormatJSON
	}

	return &ScriptStateReference{
		File:     file,
		RunEvery: runEvery,
		Format:   format,
	}, nil
}

func (s Scripts) StateReferences() ([]ScriptStateReference, error) {
	references := make([]ScriptStateReference, 0, len(s))
	for _, script := range s {
		if script == nil {
			continue
		}

		reference, err := script.stateReference()
		if err != nil {
			return nil, err
		}
		if reference == nil {
			continue
		}

		references = append(references, *reference)
	}

	return references, nil
}

func (s Scripts) PrimeState(now time.Time) int {
	checks := 0
	for _, script := range s {
		if script == nil || script.ignore() {
			continue
		}

		checked := script.prepareState(now)
		if checked {
			checks++
		}
	}

	return checks
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
		if script.ignore() {
			script.clearIgnoreFreeze()
			continue
		}

		checked := script.stateChecked
		shouldRun := script.stateShouldRun
		if !script.statePrepared {
			checked = script.prepareState(time.Now())
			shouldRun = script.stateShouldRun
		}
		if checked && !shouldRun {
			advanceAutoProgress(script.effectiveWeight())
			script.clearPreparedState()
			script.clearIgnoreFreeze()
			continue
		}

		scriptBlock := script.String()
		if len(scriptBlock) == 0 {
			advanceAutoProgress(script.effectiveWeight())
			script.clearPreparedState()
			script.clearIgnoreFreeze()
			continue
		}

		if first && dotFileHasRenderableContent() {
			if !dotFileEndsWithNewline() {
				writeRenderOutput("\n")
			}
			writeRenderOutput("\n")
		} else if !dotFileEndsWithNewline() {
			writeRenderOutput("\n")
		}
		writeRenderOutput(scriptBlock)

		if script.stateChecked {
			_ = aliaeState.WriteLastRun(script.statePath, time.Now(), script.stateFormat)
		}

		first = false
		advanceAutoProgress(script.effectiveWeight())
		script.clearPreparedState()
		script.clearIgnoreFreeze()
	}
}

func (s *Script) FreezeIgnore() bool {
	s.ifFrozen = true
	s.ifIgnoreFrozen = s.If.Ignore()
	return s.ifIgnoreFrozen
}

func (s *Script) ignore() bool {
	if s.ifFrozen {
		return s.ifIgnoreFrozen
	}

	return s.If.Ignore()
}

func (s *Script) clearIgnoreFreeze() {
	s.ifFrozen = false
	s.ifIgnoreFrozen = false
}

func (s *Script) prepareState(now time.Time) bool {
	s.statePrepared = true
	s.stateChecked = false
	s.stateShouldRun = true
	s.statePath = ""
	s.stateFormat = ""

	reference, err := s.stateReference()
	if err != nil {
		s.stateShouldRun = false
		s.stateChecked = true
		return true
	}
	if reference == nil {
		return false
	}

	statePath := aliaeState.Path(reference.File)
	shouldRun, _, checkErr := aliaeState.ShouldRun(statePath, reference.RunEvery, now)
	s.stateChecked = true
	s.statePath = statePath
	s.stateFormat = reference.Format
	if checkErr != nil {
		s.stateShouldRun = false
		return true
	}

	s.stateShouldRun = shouldRun
	return true
}

func (s *Script) clearPreparedState() {
	s.statePrepared = false
	s.stateChecked = false
	s.stateShouldRun = false
	s.statePath = ""
	s.stateFormat = ""
}

func inlineInterpreterScript(executable, switchName, script string) string {
	formatted, ok := formatString(Template(script)).(string)
	if !ok {
		return ""
	}

	command := fmt.Sprintf("%s %s %s", executable, switchName, formatted)

	runtime := currentRuntime()
	if runtime != nil && runtime.Shell == CMD {
		luaFormatted, ok := formatString(command).(string)
		if !ok {
			return ""
		}
		return fmt.Sprintf("os.execute(%s)", luaFormatted)
	}

	return command
}
