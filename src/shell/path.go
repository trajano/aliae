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
	Value       Template `yaml:"value"`
	If          If       `yaml:"if"`
	template    string
	Persist     bool `yaml:"persist"`
	Force       bool `yaml:"force"`
	IfExists    bool `yaml:"ifExists"`
	ifEvaluated bool `yaml:"-"`
	ifIgnored   bool `yaml:"-"`
}

func (p *Path) string() string {
	return formatStrategy().FormatPath(p)
}

func (p *Path) render() string {
	p.Value = p.Value.Parse()
	runtime := currentRuntime()
	if runtime == nil {
		return ""
	}
	if runtime.Path == nil {
		runtime.Path = &context.Path{}
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

		if runtime.Path.Contains(line) && !p.Force {
			continue
		}

		runtime.Path.Append(line)

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

func (p *Path) Ignore() bool {
	if p.ifEvaluated {
		return p.ifIgnored
	}

	p.ifIgnored = p.If.Ignore()
	p.ifEvaluated = true
	return p.ifIgnored
}

func pathEntryExists(entry string) bool {
	runtime := currentRuntime()

	resolved := entry
	if strings.HasPrefix(resolved, "~") {
		rem := resolved[1:]
		if len(rem) == 0 || rem[0] == '/' || rem[0] == '\\' {
			if runtime != nil && runtime.Home != "" {
				resolved = runtime.Home + rem
			} else {
				resolved = context.Home() + rem
			}
		}
	}

	if !filepath.IsAbs(resolved) && !strings.HasPrefix(resolved, "/") {
		if runtime != nil && runtime.Home != "" {
			resolved = filepath.Join(runtime.Home, resolved)
		} else {
			resolved = filepath.Join(context.Home(), resolved)
		}
	}

	if runtime != nil && runtime.OS == context.WINDOWS && isMSYS2Environment() {
		if windowsPath, ok := msysToWindowsPath(resolved); ok {
			resolved = windowsPath
		}
	}
	if runtime != nil && runtime.WSL && runtime.Shell == BASH {
		resolved = context.CleanPath(resolved)
	}

	return pathExists(resolved).exists
}

func normalizePathEntry(entry string) string {
	runtime := currentRuntime()
	if runtime == nil {
		return entry
	}

	if runtime.WSL && runtime.Shell == BASH {
		return context.CleanPath(entry)
	}

	if runtime.OS != context.WINDOWS {
		return entry
	}

	if runtime.Shell == BASH && isMSYS2Environment() {
		normalized, ok := windowsToMSYSPath(entry)
		if !ok {
			return entry
		}

		return normalized
	}

	switch runtime.Shell {
	case PWSH, POWERSHELL, CMD, NU:
		if windowsPath, ok := msysToWindowsPath(entry); ok {
			return windowsPath
		}

		return strings.ReplaceAll(entry, "/", `\`)
	default:
		return entry
	}
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
		if entry.Ignore() {
			continue
		}

		script := entry.string()
		if len(script) == 0 {
			advanceAutoProgress(1)
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
		writeRenderOutput(script)

		first = false
		advanceAutoProgress(1)
	}
}
