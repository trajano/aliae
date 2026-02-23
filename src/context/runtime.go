package context

import (
	"os"
	"runtime"
)

// used for caching runtime information
// and testing purposes
var Current *Runtime

var runtimeGOOS = runtime.GOOS

type Runtime struct {
	Path       *Path
	CDPath     *Path
	Shell      string
	OS         string
	Hostname   string
	Home       string
	Arch       string
	ConfigPath string
	ConfigDir  string
	WSL        bool
}

func Init(shell string) {
	home := Home()
	hostname, _ := os.Hostname()

	Current = &Runtime{
		Shell:    shell,
		OS:       runtimeGOOS,
		WSL:      isWSL(),
		Arch:     runtime.GOARCH,
		Home:     home,
		Hostname: hostname,
	}

	Current.Path = getPath()
	Current.CDPath = getCDPath()
}

func Home() string {
	if Current != nil {
		return Current.Home
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

func isWSL() bool {
	if runtimeGOOS != LINUX {
		return false
	}

	return os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != ""
}
