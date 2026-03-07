package appcmd

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

type stateEntry struct {
	file string
	// lastRun is preformatted for display and empty when state file is missing.
	lastRun  string
	format   aliaeState.FileFormat
	runEvery time.Duration
}

type StateListCommand struct {
	Out        io.Writer
	ConfigPath string
}

func (c StateListCommand) Execute() error {
	entries, err := referencedStateEntries(c.ConfigPath)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		_, _ = fmt.Fprintln(c.Out, "No state entries referenced in config.")
		return nil
	}

	w := tabwriter.NewWriter(c.Out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "FILE\tLAST RUN\tRUN EVERY\tFORMAT")
	for _, entry := range entries {
		lastRun := "-"
		if len(entry.lastRun) > 0 {
			lastRun = entry.lastRun
		}

		runEvery := "once"
		if entry.runEvery > 0 {
			runEvery = entry.runEvery.String()
		}

		format := string(entry.format)
		if len(format) == 0 {
			format = string(aliaeState.FormatJSON)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", entry.file, lastRun, runEvery, format)
	}
	return w.Flush()
}

type StateClearCommand struct {
	Out        io.Writer
	ConfigPath string
}

func (c StateClearCommand) Execute() error {
	entries, err := referencedStateEntries(c.ConfigPath)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		_, _ = fmt.Fprintln(c.Out, "No state entries referenced in config.")
		return nil
	}

	removed := make([]string, 0, len(entries))
	for _, entry := range entries {
		if err := os.Remove(aliaeState.Path(entry.file)); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		removed = append(removed, entry.file)
	}

	if len(removed) == 0 {
		_, _ = fmt.Fprintln(c.Out, "No state files removed.")
		return nil
	}

	slices.Sort(removed)
	for _, file := range removed {
		fmt.Fprintf(c.Out, "removed %s\n", file)
	}

	return nil
}

func referencedStateEntries(configPath string) ([]stateEntry, error) {
	runtime := context.GetCurrent()
	if runtime == nil {
		context.Init(shell.Name())
		runtime = context.GetCurrent()
	}
	if runtime == nil {
		return nil, fmt.Errorf("runtime initialization failed")
	}
	shell.SetRuntime(runtime)

	resolvedConfigPath, resolvedConfigDir := cfg.ResolveTemplateContext(configPath)
	runtime.ConfigPath = resolvedConfigPath
	runtime.ConfigDir = resolvedConfigDir

	aliae, err := cfg.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	references, err := aliae.Scripts.StateReferences()
	if err != nil {
		return nil, err
	}

	entries := make([]stateEntry, 0, len(references))
	for _, reference := range references {
		file := strings.TrimSpace(reference.File)
		if len(file) == 0 {
			continue
		}

		path := aliaeState.Path(file)
		lastRun, err := aliaeState.ReadLastRun(path)
		if err != nil {
			return nil, err
		}
		lastRunText := ""
		if lastRun != nil {
			lastRunText = lastRun.UTC().Format(time.RFC3339)
		}

		entries = append(entries, stateEntry{
			file:     file,
			runEvery: reference.RunEvery,
			format:   reference.Format,
			lastRun:  lastRunText,
		})
	}

	slices.SortFunc(entries, func(a, b stateEntry) int {
		return strings.Compare(a.file, b.file)
	})
	return entries, nil
}
