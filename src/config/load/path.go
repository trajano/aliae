package load

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

func SetTemplateConfigContext(configPath string) {
	current := context.GetCurrent()
	if current == nil {
		return
	}

	current.ConfigPath = configPath
	current.ConfigDir = ResolveConfigDir(configPath)
}

func ResolveTemplateContext(configPath string) (string, string) {
	resolvedPath := ResolveConfigPath(configPath)
	return resolvedPath, ResolveConfigDir(resolvedPath)
}

func ResolveConfigPath(configPath string) string {
	if len(configPath) == 0 {
		configPath = os.Getenv("ALIAE_CONFIG")
	}

	if len(configPath) == 0 {
		configPath = path.Join(Home(), ".aliae.yaml")
	}

	return replaceTildePrefixWithHomeDir(configPath)
}

func ResolveConfigDir(configPath string) string {
	if !strings.HasPrefix(configPath, "http://") && !strings.HasPrefix(configPath, "https://") {
		return filepath.Dir(configPath)
	}

	parsed, err := url.Parse(configPath)
	if err != nil {
		return filepath.Dir(configPath)
	}

	parsed.Path = path.Dir(parsed.Path)
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String()
}

func Home() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}

	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}

	return home
}

func replaceTildePrefixWithHomeDir(dir string) string {
	if !strings.HasPrefix(dir, "~") {
		return dir
	}

	rem := dir[1:]
	if len(rem) == 0 || isSeparator(rem[0]) {
		return Home() + rem
	}

	return dir
}

func isSeparator(c uint8) bool {
	if c == '/' {
		return true
	}

	if runtime.GOOS == context.WINDOWS && c == '\\' {
		return true
	}

	return false
}
