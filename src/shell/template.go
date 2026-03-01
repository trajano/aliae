package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Template string

var (
	hasCommandCache   = map[string]bool{}
	hasCommandCacheMu sync.RWMutex

	pathExistsCache   = map[string]pathInfo{}
	pathExistsCacheMu sync.RWMutex
)

type pathInfo struct {
	exists bool
	isDir  bool
}

func (t Template) Parse() Template {
	value, err := parse(string(t), context.Current)
	if err != nil {
		return t
	}

	return Template(value)
}

func (t Template) String() string {
	return string(t.Parse())
}

func parse(text string, ctx any) (string, error) {
	if !strings.Contains(text, "{{") || !strings.Contains(text, "}}") {
		return text, nil
	}

	parsedTemplate, err := template.New("alias").Option("missingkey=zero").Funcs(funcMap()).Parse(text)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = parsedTemplate.Execute(buffer, ctx)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"isPwshOption":      isPwshOption,
		"isPwshScope":       isPwshScope,
		"formatString":      formatString,
		"formatArray":       formatArray,
		"escapeString":      escapeString,
		"env":               os.Getenv,
		"match":             match,
		"hasCommand":        hasCommand,
		"hasCommandNoCache": hasCommandNoCache,
		"fileExists":        fileExists,
		"dirExists":         dirExists,
		"isDir":             isDir,
		"progress":          progress,
	}
	return funcMap
}

func formatString(variable any) any {
	switch variable.(type) {
	case string, Template, Option:
		return fmt.Sprintf(`"%s"`, escapeString(variable))
	default:
		return variable
	}
}

func splitString(variable any) any {
	switch variable := variable.(type) {
	case string:
		variable = strings.TrimSpace(variable)
		if len(variable) == 0 {
			return []string{variable}
		}

		if strings.Contains(variable, "\n") {
			return strings.Split(variable, "\n")
		}

		return strings.Fields(variable)
	case Template:
		return splitString(variable.String())
	default:
		return variable
	}
}

func formatArray(variable any, delim ...string) any {
	delimiter := " "
	if len(delim) > 0 {
		delimiter = delim[0]
	}

	switch variable := variable.(type) {
	case string:
		split := splitString(variable).([]string)
		array := []string{}

		for _, value := range split {
			array = append(array, formatString(value).(string))
		}

		return strings.Join(array, delimiter)
	case Template:
		return formatArray(variable.String())
	default:
		return variable
	}
}

func escapeString(variable any) any {
	clean := func(v string) string {
		switch context.Current.Shell {
		case PWSH, POWERSHELL:
			return strings.NewReplacer(
				"`", "``",
				`"`, "`\"",
			).Replace(v)
		default:
			return strings.NewReplacer(
				`\`, `\\`,
				`"`, `\"`,
			).Replace(v)
		}
	}

	switch v := variable.(type) {
	case Template:
		value := v.String()
		return clean(value)
	case string:
		return clean(v)
	default:
		return variable
	}
}

func match(variable string, values ...string) bool {
	return slices.Contains(values, variable)
}

func hasCommand(command string) bool {
	hasCommandCacheMu.RLock()
	cached, ok := hasCommandCache[command]
	hasCommandCacheMu.RUnlock()
	if ok {
		return cached
	}

	result := hasCommandNoCache(command)

	hasCommandCacheMu.Lock()
	hasCommandCache[command] = result
	hasCommandCacheMu.Unlock()

	return result
}

func hasCommandNoCache(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func fileExists(path string) bool {
	info := pathExists(resolveFromHome(path))
	return info.exists && !info.isDir
}

func dirExists(path string) bool {
	info := pathExists(resolveFromHome(path))
	return info.exists && info.isDir
}

func pathExists(path string) pathInfo {
	pathExistsCacheMu.RLock()
	cached, OK := pathExistsCache[path]
	pathExistsCacheMu.RUnlock()
	if OK {
		return cached
	}

	info, err := statWithTimeout(path, statTimeout)
	result := pathInfo{
		exists: err == nil,
		isDir:  err == nil && info.IsDir(),
	}

	pathExistsCacheMu.Lock()
	pathExistsCache[path] = result
	pathExistsCacheMu.Unlock()

	return result
}

func resolveFromHome(path string) string {
	if filepath.IsAbs(path) || strings.HasPrefix(path, "/") {
		return path
	}

	return filepath.Join(context.Home(), path)
}

func clearPathExistsCache() {
	pathExistsCacheMu.Lock()
	defer pathExistsCacheMu.Unlock()
	pathExistsCache = map[string]pathInfo{}
}

func clearHasCommandCache() {
	hasCommandCacheMu.Lock()
	defer hasCommandCacheMu.Unlock()
	hasCommandCache = map[string]bool{}
}

// ResetTemplateCaches clears in-memory caches used by template helper functions.
func ResetTemplateCaches() {
	clearHasCommandCache()
	clearPathExistsCache()
}

func executableExtension() string {
	if runtime.GOOS == context.WINDOWS {
		return ".cmd"
	}

	return ""
}

func isDir(path string) bool {
	info, err := statWithTimeout(path, statTimeout)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func progress(value any) string {
	state := 1
	percentage := 0

	switch v := value.(type) {
	case int:
		percentage = v
	case int8:
		percentage = int(v)
	case int16:
		percentage = int(v)
	case int32:
		percentage = int(v)
	case int64:
		percentage = int(v)
	case uint:
		percentage = int(v)
	case uint8:
		percentage = int(v)
	case uint16:
		percentage = int(v)
	case uint32:
		percentage = int(v)
	case uint64:
		percentage = int(v)
	case string:
		text := strings.TrimSpace(v)
		if strings.EqualFold(text, "reset") {
			state = 0
			percentage = 0
			break
		}
		parsed, err := strconv.Atoi(text)
		if err != nil {
			return ""
		}
		percentage = parsed
	default:
		return ""
	}

	if percentage < 0 || percentage > 100 {
		return ""
	}

	switch context.Current.Shell {
	case PWSH, POWERSHELL:
		return fmt.Sprintf(`[Console]::Out.Write("$([char]27)]9;4;%d;%d$([char]7)")`, state, percentage)
	default:
		return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
	}
}
