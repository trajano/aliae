package context

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Path []string

func getPath() *Path {
	if Current != nil {
		return Current.Path
	}

	path := &Path{}
	paths := os.Getenv("PATH")

	for _, p := range strings.Split(paths, PathDelimiter()) {
		path.Append(cleanPath(p))
	}

	return path
}

func (p *Path) Append(path string) {
	clean := cleanPath(path)
	if len(clean) == 0 || p.Contains(clean) {
		return
	}

	current := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s%s%s", clean, PathDelimiter(), current))

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
		if normalized, OK := windowsToMSYSPath(path); OK {
			path = normalized
		}
	}

	return strings.TrimRight(path, `/\`)
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

func isASCIIAlpha(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}
