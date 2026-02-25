package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/registry"
)

type Paths []*Path

type Path struct {
	Value    Template `yaml:"value"`
	If       If       `yaml:"if"`
	template string
	Persist  bool `yaml:"persist"`
	Force    bool `yaml:"force"`
	IfExists bool `yaml:"ifExists"`
}

func (p *Path) string() string {
	switch context.Current.Shell {
	case ZSH:
		return p.zsh().render()
	case BASH:
		return p.bash().render()
	case PWSH, POWERSHELL:
		return p.pwsh().render()
	case NU:
		return p.nu().render()
	case FISH:
		return p.fish().render()
	case TCSH:
		return p.tcsh().render()
	case XONSH:
		return p.xonsh().render()
	case CMD:
		return p.cmd().render()
	default:
		return ""
	}
}

func (p *Path) render() string {
	p.Value = p.Value.Parse()

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

		if context.Current.Path.Contains(line) && !p.Force {
			continue
		}

		context.Current.Path.Append(line)

		if !first {
			builder.WriteString("\n")
		}

		if p.Persist {
			registry.PersistPathEntry(rawLine)
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

func pathEntryExists(entry string) bool {
	resolved := entry
	if strings.HasPrefix(resolved, "~") {
		rem := resolved[1:]
		if len(rem) == 0 || rem[0] == '/' || rem[0] == '\\' {
			resolved = context.Home() + rem
		}
	}

	if !filepath.IsAbs(resolved) && !strings.HasPrefix(resolved, "/") {
		resolved = filepath.Join(context.Home(), resolved)
	}

	if context.Current != nil && context.Current.OS == context.WINDOWS && isMSYS2Environment() {
		if windowsPath, ok := msysToWindowsPath(resolved); ok {
			resolved = windowsPath
		}
	}

	return pathExists(resolved).exists
}

func normalizePathEntry(entry string) string {
	if context.Current == nil || context.Current.OS != context.WINDOWS || context.Current.Shell != BASH || !isMSYS2Environment() {
		return entry
	}

	normalized, ok := windowsToMSYSPath(entry)
	if !ok {
		return entry
	}

	return normalized
}

func isMSYS2Environment() bool {
	return os.Getenv("MSYSTEM") != ""
}

func windowsToMSYSPath(path string) (string, bool) {
	if len(path) < 3 || path[1] != ':' {
		return "", false
	}

	drive := path[0]
	if !isASCIIAlpha(drive) {
		return "", false
	}

	if path[2] != '\\' && path[2] != '/' {
		return "", false
	}

	rest := strings.ReplaceAll(path[2:], `\`, `/`)
	return fmt.Sprintf("/%s%s", strings.ToLower(string(drive)), rest), true
}

func msysToWindowsPath(path string) (string, bool) {
	if len(path) < 4 || path[0] != '/' || path[2] != '/' {
		return "", false
	}

	drive := path[1]
	if !isASCIIAlpha(drive) {
		return "", false
	}

	rest := strings.ReplaceAll(path[3:], `/`, `\`)
	return fmt.Sprintf("%s:\\%s", strings.ToUpper(string(drive)), rest), true
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}

func (p Paths) Render() {
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
