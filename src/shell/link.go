package shell

import (
	"os"
	"path/filepath"
)

type Links []*Link

type Link struct {
	Name     Template `yaml:"name"`
	Target   Template `yaml:"target"`
	If       If       `yaml:"if"`
	template string
	MkDir    bool `yaml:"mkdir"`
	force    bool
}

func (l *Link) string() string {
	// avoid parsing multiple times
	l.Name = l.Name.Parse()
	l.Target = l.Target.Parse()

	// do not process if the link already exists or the target does not exist
	if l.exists(string(l.Name)) || (!l.force && !l.exists(string(l.Target))) {
		return ""
	}

	if l.MkDir {
		l.buildPath()
	}

	return renderStrategy().RenderLink(l)
}

func (l *Link) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (l *Link) buildPath() {
	parent := filepath.Dir(string(l.Name))

	_, err := os.Stat(parent)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		_ = os.MkdirAll(parent, 0644)
	}
}

func (l *Link) render() string {
	script, err := parse(l.template, l)
	if err != nil {
		return err.Error()
	}

	return script
}

func (l Links) Render() {
	if len(l) == 0 {
		return
	}

	first := true
	for _, link := range l {
		if link.If.Ignore() {
			continue
		}
		script := link.string()
		if len(script) == 0 {
			advanceAutoProgress(1)
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
		DotFile.WriteString(script)

		first = false
		advanceAutoProgress(1)
	}
}
