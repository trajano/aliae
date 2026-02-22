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
	switch context.Current.Shell {
	case ZSH:
		return p.zsh().render()
	case BASH:
		return p.bash().render()
	case FISH:
		return p.fish().render()
	case TCSH:
		return p.tcsh().render()
	case XONSH:
		return p.xonsh().render()
	default:
		return ""
	}
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

func (p CDPaths) Render() {
	if len(p) == 0 {
		return
	}

	first := true
	for _, entry := range p {
		if entry.If.Ignore() {
			continue
		}

		script := entry.string()
		if len(script) == 0 {
			continue
		}

		if first && DotFile.Len() > 0 {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString("\n")
		DotFile.WriteString(script)

		first = false
	}
}
