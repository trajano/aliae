package context

import (
	"os"
	"runtime"
)

// used for caching runtime information
// and testing purposes
var Current *Runtime

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
}

func Init(shell string) {
	home := Home()
	hostname, _ := os.Hostname()

	Current = &Runtime{
		Shell:    shell,
		OS:       runtime.GOOS,
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
