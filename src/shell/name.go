package shell

import (
	"fmt"
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
	name, _ = resolveShellName(executable, name, func() (string, error) {
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

func NameVerbose() (string, []string) {
	pid := os.Getppid()
	p, _ := process.NewProcess(int32(pid))
	name, err := p.Name()
	if err != nil {
		return UNKNOWN, []string{fmt.Sprintf("failed to resolve parent process name: %v", err)}
	}

	executable, _ := os.Executable()
	executable = filepath.Base(executable)
	trace := []string{
		fmt.Sprintf("executable=%s", executable),
		fmt.Sprintf("layer 0 pid=%d name=%s", p.Pid, name),
	}

	layer := 0
	resolved, resolutionTrace := resolveShellName(executable, name, func() (string, error) {
		layer++
		p, _ = p.Parent()
		if p == nil {
			trace = append(trace, fmt.Sprintf("layer %d pid=? name=<none>", layer))
			return "", nil
		}

		parentName, parentErr := p.Name()
		if parentErr != nil {
			trace = append(trace, fmt.Sprintf("layer %d pid=%d error=%v", layer, p.Pid, parentErr))
			return "", parentErr
		}

		trace = append(trace, fmt.Sprintf("layer %d pid=%d name=%s", layer, p.Pid, parentName))
		return parentName, nil
	})

	trace = append(trace, resolutionTrace...)

	if resolved == "" {
		trace = append(trace, "resolved=unknown")
		return UNKNOWN, trace
	}

	trimmed := strings.TrimSuffix(resolved, ".exe")
	trace = append(trace, fmt.Sprintf("resolved=%s", trimmed))
	return trimmed, trace
}

func shouldUseParentShell(name, executable string) bool {
	normalizedName := strings.TrimSuffix(filepath.Base(name), ".exe")
	normalizedExecutable := strings.TrimSuffix(filepath.Base(executable), ".exe")
	return normalizedName == normalizedExecutable
}

func resolveShellName(executable, current string, next func() (string, error)) (string, []string) {
	name := current
	trace := []string{}
	for shouldUseParentShell(name, executable) {
		trace = append(trace, fmt.Sprintf("shim-detected name=%s", name))
		parentName, err := next()
		if err != nil || parentName == "" {
			trace = append(trace, "stopped: no further parent shell available")
			return "", trace
		}
		name = parentName
	}
	return name, trace
}
