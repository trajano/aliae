package context

import (
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

// used for caching runtime information
// and testing purposes
var currentRuntime atomic.Pointer[Runtime]

var runtimeGOOS = runtime.GOOS

type Runtime struct {
	Path       *Path
	CDPath     *Path
	Env        map[string]string
	Var        map[string]any
	Shell      string
	OS         string
	Hostname   string
	Home       string
	Arch       string
	ConfigPath string
	ConfigDir  string
	Cygpath    string
	ShellLike  bool
	WSL        bool
}

func Init(shell string) {
	SetCurrent(NewRuntime(shell))
}

func NewRuntime(shell string) *Runtime {
	home := Home()
	hostname, _ := os.Hostname()

	current := &Runtime{
		Shell:     shell,
		ShellLike: isShellLike(shell),
		OS:        runtimeGOOS,
		WSL:       isWSL(),
		Arch:      runtime.GOARCH,
		Home:      home,
		Hostname:  hostname,
		Env:       getEnvironment(),
		Var:       map[string]any{},
		Cygpath:   CygpathInternal,
	}

	current.Path = getPath()
	current.CDPath = getCDPath()
	return current
}

func isShellLike(shell string) bool {
	switch strings.ToLower(strings.TrimSpace(shell)) {
	case "bash", "zsh", "fish", "tcsh", "pwsh", "powershell":
		return true
	default:
		return false
	}
}

func Home() string {
	if current := GetCurrent(); current != nil {
		return current.Home
	}

	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}
	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

func SetCurrent(current *Runtime) {
	currentRuntime.Store(current)
}

func GetCurrent() *Runtime {
	return currentRuntime.Load()
}

func isWSL() bool {
	if runtimeGOOS != LINUX {
		return false
	}

	return os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != ""
}

func getEnvironment() map[string]string {
	env := make(map[string]string)

	for _, value := range os.Environ() {
		key, val, found := strings.Cut(value, "=")
		if !found {
			continue
		}
		env[key] = val
	}

	return env
}
