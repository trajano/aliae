package shell

type Scripts []*Script

type Script struct {
	Value  Template `yaml:"value"`
	If     If       `yaml:"if"`
	Weight *float64 `yaml:"weight"`
}

func (s *Script) String() string {
	script := s.Value.Parse()
	return string(script)
}

func (s *Script) effectiveWeight() float64 {
	if s.Weight == nil {
		return 1
	}

	if *s.Weight < 1 {
		return 1
	}

	return *s.Weight
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
