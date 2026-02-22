package context

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Path []string

var runCygpath = func(path string) (string, error) {
	output, err := exec.Command("cygpath", "-u", path).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getPath() *Path {
	return getEnvPath("PATH")
}

func getCDPath() *Path {
	return getEnvPath("CDPATH")
}

func getEnvPath(envName string) *Path {
	path := &Path{}
	paths := os.Getenv(envName)

	for _, p := range splitPathEntries(paths) {
		clean := cleanPath(p)
		if len(clean) == 0 || slices.Contains(*path, clean) {
			continue
		}

		*path = append(*path, clean)
	}

	return path
}

func (p *Path) Append(path string) {
	p.appendWithEnv("PATH", path, true)
}

func (p *Path) AppendCDPath(path string) {
	p.appendWithEnv("CDPATH", path, false)
}

func (p *Path) appendWithEnv(envName, path string, prepend bool) {
	clean := cleanPath(path)
	if len(clean) == 0 || p.Contains(clean) {
		return
	}

	current := os.Getenv(envName)
	switch {
	case current == "":
		os.Setenv(envName, clean)
	case prepend:
		os.Setenv(envName, fmt.Sprintf("%s%s%s", clean, PathDelimiter(), current))
	default:
		os.Setenv(envName, fmt.Sprintf("%s%s%s", current, PathDelimiter(), clean))
	}

	*p = append(*p, clean)
}

func (p *Path) Contains(path string) bool {
	return slices.Contains(*p, cleanPath(path))
}

func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return path
	}

	if isMSYS2Shell() {
		if normalized, err := runCygpath(path); err == nil && normalized != "" {
			path = normalized
		}
	}

	return strings.TrimRight(path, `/\`)
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}

func splitPathEntries(paths string) []string {
	if Current == nil || Current.OS != WINDOWS {
		return strings.Split(paths, PathDelimiter())
	}

	result := make([]string, 0)
	start := 0

	for i := range len(paths) {
		switch paths[i] {
		case ';':
			result = append(result, paths[start:i])
			start = i + 1
		case ':':
			if i == start+1 && start < len(paths) && isASCIIAlpha(paths[start]) {
				// Keep Windows drive-letter paths together (for example C:\...).
				continue
			}
			result = append(result, paths[start:i])
			start = i + 1
		}
	}

	result = append(result, paths[start:])
	return result
}
