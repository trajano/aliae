package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

type CDPaths []*CDPath

type CDPath struct {
	Value    Template `yaml:"value"`
	If       If       `yaml:"if"`
	template string
	Force    bool `yaml:"force"`
	IfExists bool `yaml:"ifExists"`
}

func (p *CDPath) string() string {
	return newShellFactory().strategy().renderCDPath(p)
}

func (p *CDPath) render() string {
	p.Value = p.Value.Parse()
	if context.Current.CDPath == nil {
		context.Current.CDPath = &context.Path{}
	}

	var builder strings.Builder
	ctx := struct {
		Value string
	}{}

	splitted := strings.Split(string(p.Value), "\n")

	first := true
	for _, line := range splitted {
		rawLine := strings.TrimSpace(line)
		if len(rawLine) == 0 {
			continue
		}

		if p.IfExists && !pathEntryExists(rawLine) {
			continue
		}

		line = normalizePathEntry(rawLine)

		if context.Current.CDPath.Contains(line) && !p.Force {
			continue
		}

		context.Current.CDPath.AppendCDPath(line)

		if !first {
			builder.WriteString("\n")
		}

		ctx.Value = line
		script, err := parse(p.template, ctx)
		if err != nil {
			builder.WriteString(err.Error())
		}

		builder.WriteString(script)
		first = false
	}

	return builder.String()
}

func cdpathCurrentDirScript() string {
	return newShellFactory().strategy().renderCDPathCurrentDirScript()
}

func (p CDPaths) Render() {
	if len(p) == 0 {
		return
	}

	first := true
	rendered := false
	for _, entry := range p {
		if entry.If.Ignore() {
			continue
		}

		script := entry.string()
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
		rendered = true
		advanceAutoProgress(1)
	}

	// Some shells stop treating the current directory as an implicit fallback
	// when CDPATH/cdpath is set and "." is missing.
	if rendered && !context.Current.CDPath.Contains(".") {
		if !dotFileEndsWithNewline() {
			DotFile.WriteString("\n")
		}
		DotFile.WriteString(cdpathCurrentDirScript())
		context.Current.CDPath.AppendCDPath(".")
	}
}
