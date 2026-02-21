package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/process"
)

const (
	UNKNOWN = "unknown"
)

func Name() string {
	pid := os.Getppid()
	p, _ := process.NewProcess(int32(pid))
	name, err := p.Name()
	if err != nil {
		return UNKNOWN
	}

	executable, _ := os.Executable()
	executable = filepath.Base(executable)
	name = resolveShellName(executable, name, func() (string, error) {
		p, _ = p.Parent()
		if p == nil {
			return "", nil
		}
		return p.Name()
	})

	if name == "" {
		return UNKNOWN
	}

	if err != nil {
		return UNKNOWN
	}

	return strings.TrimSuffix(name, ".exe")
}

func shouldUseParentShell(name, executable string) bool {
	normalizedName := strings.TrimSuffix(filepath.Base(name), ".exe")
	normalizedExecutable := strings.TrimSuffix(filepath.Base(executable), ".exe")
	return normalizedName == normalizedExecutable
}

func resolveShellName(executable, current string, next func() (string, error)) string {
	name := current
	for shouldUseParentShell(name, executable) {
		parentName, err := next()
		if err != nil || parentName == "" {
			return ""
		}
		name = parentName
	}
	return name
}
