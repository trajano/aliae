package context

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	internalcygpath "github.com/jandedobbeleer/aliae/src/cygpath"
)

type Path []string

var runExternalCygpath = func(path string) (string, error) {
	output, err := exec.Command("cygpath", "-u", path).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

var cleanPathCache sync.Map

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

	cacheKey := cleanPathCacheKey(path)
	if cached, ok := cleanPathCache.Load(cacheKey); ok {
		return cached.(string)
	}

	if isMSYS2Shell() {
		switch cygpathMode() {
		case CygpathExternal:
			if normalized, err := runExternalCygpath(path); err == nil && normalized != "" {
				path = normalized
			}
		default:
			path = internalcygpath.ToUnix(path)
		}
	}

	clean := strings.TrimRight(path, `/\`)
	cleanPathCache.Store(cacheKey, clean)
	return clean
}

func cygpathMode() string {
	current := GetCurrent()
	if current == nil {
		return CygpathInternal
	}

	mode := NormalizeCygpathMode(current.Cygpath)
	if mode == CygpathExternal {
		return CygpathExternal
	}

	return CygpathInternal
}

func cleanPathCacheKey(path string) string {
	current := GetCurrent()
	if current == nil {
		return "|:" + os.Getenv("MSYSTEM") + ":" + path
	}

	return current.OS + "|" + current.Shell + "|" + cygpathMode() + ":" + os.Getenv("MSYSTEM") + ":" + path
}

func clearCleanPathCache() {
	cleanPathCache.Range(func(key, _ any) bool {
		cleanPathCache.Delete(key)
		return true
	})
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}

func splitPathEntries(paths string) []string {
	current := GetCurrent()
	if current == nil || current.OS != WINDOWS {
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
