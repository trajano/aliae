package appcmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	initpkg "github.com/jandedobbeleer/aliae/src/init"
	"github.com/jandedobbeleer/aliae/src/shell"
	"golang.org/x/term"
)

type GetShellCommand struct {
	Out io.Writer
}

func (c GetShellCommand) Execute() error {
	_, err := fmt.Fprintln(c.Out, shell.Name())
	return err
}

type GetResolvedConfigCommand struct {
	Out        io.Writer
	ConfigPath string
}

func (c GetResolvedConfigCommand) Execute() error {
	output, err := RenderResolvedConfigYAML(c.ConfigPath)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.Out, output)
	return err
}

type GetVariablesCommand struct {
	Out        io.Writer
	ConfigPath string
}

func (c GetVariablesCommand) Execute() error {
	return printVariableDiagnostics(c.Out, c.ConfigPath)
}

type BenchmarkCommand struct {
	Out            io.Writer
	ConfigPath     string
	BenchmarkShell string
	NoCache        bool
}

const (
	yamlKeyForce    = "force"
	yamlKeyIfExists = "ifExists"
	yamlKeyName     = "name"
	yamlKeyType     = "type"
	yamlKeyValue    = "value"
)

func (c BenchmarkCommand) Execute() error {
	return printBenchmark(c.Out, c.ConfigPath, c.BenchmarkShell, c.NoCache)
}

func RenderResolvedConfigYAML(configPath string) (string, error) {
	aliae, err := cfg.LoadConfig(configPath)
	if err != nil {
		return "", err
	}

	output := yaml.MapSlice{}
	if aliases := toAliasOutput(aliae.Aliae); len(aliases) > 0 {
		output = append(output, yaml.MapItem{Key: "alias", Value: aliases})
	}
	if vars := toVarOutput(aliae.Vars); len(vars) > 0 {
		output = append(output, yaml.MapItem{Key: "var", Value: vars})
	}
	if envs := toEnvOutput(aliae.Envs); len(envs) > 0 {
		output = append(output, yaml.MapItem{Key: "env", Value: envs})
	}
	if paths := toPathOutput(aliae.Paths); len(paths) > 0 {
		output = append(output, yaml.MapItem{Key: "path", Value: paths})
	}
	if cdpaths := toCDPathOutput(aliae.CDPaths); len(cdpaths) > 0 {
		output = append(output, yaml.MapItem{Key: "cdpath", Value: cdpaths})
	}
	if scripts := toScriptOutput(aliae.Scripts); len(scripts) > 0 {
		output = append(output, yaml.MapItem{Key: "script", Value: scripts})
	}
	if links := toLinkOutput(aliae.Links); len(links) > 0 {
		output = append(output, yaml.MapItem{Key: "link", Value: links})
	}
	if aliae.Progress.Enabled {
		output = append(output, yaml.MapItem{Key: "progress", Value: aliae.Progress})
	}
	if aliae.StatTimeout > 0 {
		output = append(output, yaml.MapItem{Key: "stat_timeout", Value: aliae.StatTimeout.String()})
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func toAliasOutput(items shell.Aliae) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, alias := range items {
		if alias == nil {
			continue
		}

		item := yaml.MapSlice{{Key: yamlKeyName, Value: alias.Name}}
		if len(alias.Type) > 0 {
			item = append(item, yaml.MapItem{Key: yamlKeyType, Value: alias.Type})
		}
		item = append(item, yaml.MapItem{Key: yamlKeyValue, Value: alias.Value})
		if len(alias.Description) > 0 {
			item = append(item, yaml.MapItem{Key: "description", Value: alias.Description})
		}
		if alias.Force {
			item = append(item, yaml.MapItem{Key: yamlKeyForce, Value: true})
		}
		if len(alias.If) > 0 {
			item = append(item, yaml.MapItem{Key: "if", Value: alias.If})
		}
		if len(alias.Option) > 0 {
			item = append(item, yaml.MapItem{Key: "option", Value: alias.Option})
		}
		if len(alias.Scope) > 0 {
			item = append(item, yaml.MapItem{Key: "scope", Value: alias.Scope})
		}

		output = append(output, item)
	}

	return output
}

func toEnvOutput(items shell.Envs) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, env := range items {
		if env == nil {
			continue
		}

		item := yaml.MapSlice{{Key: yamlKeyName, Value: env.Name}}
		if len(env.Type) > 0 {
			item = append(item, yaml.MapItem{Key: yamlKeyType, Value: env.Type})
		}
		item = append(item, yaml.MapItem{Key: yamlKeyValue, Value: env.Value})
		if len(env.Delimiter) > 0 {
			item = append(item, yaml.MapItem{Key: "delimiter", Value: env.Delimiter})
		}
		if len(env.If) > 0 {
			item = append(item, yaml.MapItem{Key: "if", Value: env.If})
		}
		if env.IfExists {
			item = append(item, yaml.MapItem{Key: yamlKeyIfExists, Value: true})
		}
		if env.IsPath {
			item = append(item, yaml.MapItem{Key: "isPath", Value: true})
		}
		if env.Persist {
			item = append(item, yaml.MapItem{Key: "persist", Value: true})
		}

		output = append(output, item)
	}

	return output
}

func toVarOutput(items cfg.Vars) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, variable := range items {
		if variable == nil {
			continue
		}

		item := yaml.MapSlice{
			{Key: yamlKeyName, Value: variable.Name},
			{Key: yamlKeyValue, Value: variable.Value},
		}
		if len(variable.If) > 0 {
			item = append(item, yaml.MapItem{Key: "if", Value: variable.If})
		}

		output = append(output, item)
	}

	return output
}

func toPathOutput(items shell.Paths) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: yamlKeyValue, Value: item.Value}}
		if item.Force {
			entry = append(entry, yaml.MapItem{Key: yamlKeyForce, Value: true})
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.IfExists {
			entry = append(entry, yaml.MapItem{Key: yamlKeyIfExists, Value: true})
		}
		if item.Persist {
			entry = append(entry, yaml.MapItem{Key: "persist", Value: true})
		}

		output = append(output, entry)
	}

	return output
}

func toCDPathOutput(items shell.CDPaths) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: yamlKeyValue, Value: item.Value}}
		if item.Force {
			entry = append(entry, yaml.MapItem{Key: yamlKeyForce, Value: true})
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.IfExists {
			entry = append(entry, yaml.MapItem{Key: yamlKeyIfExists, Value: true})
		}

		output = append(output, entry)
	}

	return output
}

func toScriptOutput(items shell.Scripts) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, script := range items {
		if script == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: yamlKeyValue, Value: script.Value}}
		if len(script.Type) > 0 {
			entry = append(entry, yaml.MapItem{Key: yamlKeyType, Value: script.Type})
		}
		if len(script.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: script.If})
		}
		if script.Weight > 0 {
			entry = append(entry, yaml.MapItem{Key: "weight", Value: script.Weight})
		}
		if len(script.State.File) > 0 {
			stateEntry := yaml.MapSlice{{Key: "file", Value: script.State.File}}
			if len(script.State.RunEvery) > 0 {
				stateEntry = append(stateEntry, yaml.MapItem{Key: "runEvery", Value: script.State.RunEvery})
			}
			if len(script.State.Format) > 0 {
				stateEntry = append(stateEntry, yaml.MapItem{Key: "format", Value: script.State.Format})
			}
			entry = append(entry, yaml.MapItem{Key: "state", Value: stateEntry})
		}

		output = append(output, entry)
	}

	return output
}

func toLinkOutput(items shell.Links) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{
			{Key: yamlKeyName, Value: item.Name},
			{Key: "target", Value: item.Target},
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.MkDir {
			entry = append(entry, yaml.MapItem{Key: "mkdir", Value: true})
		}

		output = append(output, entry)
	}

	return output
}

// printVariableDiagnostics prints runtime and template diagnostics for 'aliae get variables'.
func printVariableDiagnostics(out io.Writer, configPath string) error {
	shellName, trace := shell.NameVerbose()
	context.Init(shellName)
	runtime := context.GetCurrent()
	if runtime == nil {
		return fmt.Errorf("runtime initialization failed")
	}
	shell.SetRuntime(runtime)

	resolvedConfigPath, configDir := cfg.ResolveTemplateContext(configPath)
	runtime.ConfigPath = resolvedConfigPath
	runtime.ConfigDir = configDir

	stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
	stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))

	fmt.Fprintln(out, "aliae get variables")
	fmt.Fprintf(out, "tty.stdin=%t\n", stdinTTY)
	fmt.Fprintf(out, "tty.stdout=%t\n", stdoutTTY)
	fmt.Fprintf(out, "template.Shell=%s\n", runtime.Shell)
	fmt.Fprintf(out, "template.ShellLike=%t\n", runtime.ShellLike)
	fmt.Fprintf(out, "template.OS=%s\n", runtime.OS)
	fmt.Fprintf(out, "template.WSL=%t\n", runtime.WSL)
	fmt.Fprintf(out, "template.Hostname=%s\n", runtime.Hostname)
	fmt.Fprintf(out, "template.Home=%s\n", runtime.Home)
	fmt.Fprintf(out, "template.Arch=%s\n", runtime.Arch)
	fmt.Fprintf(out, "template.ConfigPath=%s\n", runtime.ConfigPath)
	fmt.Fprintf(out, "template.ConfigDir=%s\n", runtime.ConfigDir)
	for _, line := range trace {
		fmt.Fprintf(out, "shell.trace=%s\n", line)
	}

	return nil
}

const (
	benchmarkKeyPrefix              = "benchmark."
	benchmarkKeyShell               = benchmarkKeyPrefix + "shell"
	benchmarkKeyCacheUsed           = benchmarkKeyPrefix + "cache_used"
	benchmarkKeyCygpath             = benchmarkKeyPrefix + "cygpath"
	benchmarkKeyValidateConfig      = benchmarkKeyPrefix + "validate_config"
	benchmarkKeyTotal               = benchmarkKeyPrefix + "total"
	benchmarkStepLoadConfig         = "load_config"
	benchmarkStepEvaluateVars       = "evaluate_vars"
	benchmarkStepRenderConfig       = "render_config"
	benchmarkStepValidateConfig     = "validate_config"
	benchmarkStepGenerateInitPrefix = "generate_init_"
	benchmarkStepInitVisitPrefix    = "init.visit."
	benchmarkStepInitVisitDuration  = ".duration"
	benchmarkStepInitVisitCount     = ".count"
)

func printBenchmark(out io.Writer, configPath, benchmarkShell string, noCache bool) error {
	cfg.SetCacheBypass(noCache)
	defer cfg.SetCacheBypass(false)

	type benchmarkStep struct {
		name     string
		duration time.Duration
		count    int
		hasCount bool
	}

	steps := make([]benchmarkStep, 0, 4)
	record := func(name string, run func() error) error {
		shell.ResetTemplateCaches()
		start := time.Now()
		if err := run(); err != nil {
			return err
		}
		steps = append(steps, benchmarkStep{name: name, duration: time.Since(start)})
		return nil
	}

	totalStart := time.Now()
	var aliae *cfg.Aliae
	cacheUsed := false
	cygpathMode := context.CygpathInternal

	if err := record(benchmarkStepLoadConfig, func() error {
		var err error
		aliae, err = cfg.LoadConfigWithoutVars(configPath)
		if err != nil {
			return err
		}
		cacheUsed = cfg.LastLoadUsedCache()
		cygpathMode = context.NormalizeCygpathMode(aliae.Cygpath)
		return nil
	}); err != nil {
		return err
	}

	if err := record(benchmarkStepEvaluateVars, func() error {
		return aliae.ComputeVars()
	}); err != nil {
		return err
	}

	if err := record(benchmarkStepRenderConfig, func() error {
		_, err := RenderResolvedConfigYAML(configPath)
		return err
	}); err != nil {
		return err
	}

	validateSkipped := true
	if noCache {
		validateSkipped = false
		if err := record(benchmarkStepValidateConfig, func() error {
			return cfg.ValidateConfig(configPath)
		}); err != nil {
			return err
		}
	}

	if len(benchmarkShell) > 0 {
		stepName := benchmarkStepGenerateInitPrefix + benchmarkShell
		if err := record(stepName, func() error {
			initpkg.SetProgressWriter(io.Discard)
			defer initpkg.SetProgressWriter(os.Stderr)

			initSteps, initVisits, err := initpkg.Benchmark(configPath, benchmarkShell)
			if err != nil {
				return fmt.Errorf("init benchmark failed for shell %s: %w", benchmarkShell, err)
			}

			for _, step := range initSteps {
				steps = append(steps, benchmarkStep{
					name:     step.Name,
					duration: step.Duration,
				})
			}
			for _, visit := range initVisits {
				section := string(visit.Section)
				steps = append(steps,
					benchmarkStep{
						name:     benchmarkStepInitVisitPrefix + section + benchmarkStepInitVisitDuration,
						duration: visit.Duration,
					},
					benchmarkStep{
						name:     benchmarkStepInitVisitPrefix + section + benchmarkStepInitVisitCount,
						count:    visit.Count,
						hasCount: true,
					},
				)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	total := time.Since(totalStart)

	fmt.Fprintln(out, "aliae get benchmark")
	if len(benchmarkShell) > 0 {
		fmt.Fprintf(out, "%s=%s\n", benchmarkKeyShell, benchmarkShell)
	}
	fmt.Fprintf(out, "%s=%t\n", benchmarkKeyCacheUsed, cacheUsed)
	fmt.Fprintf(out, "%s=%s\n", benchmarkKeyCygpath, cygpathMode)
	if validateSkipped {
		fmt.Fprintf(out, "%s=skipped\n", benchmarkKeyValidateConfig)
	}
	for _, step := range steps {
		if step.hasCount {
			fmt.Fprintf(out, "%s%s=%d\n", benchmarkKeyPrefix, step.name, step.count)
			continue
		}
		fmt.Fprintf(out, "%s%s=%s\n", benchmarkKeyPrefix, step.name, step.duration)
	}
	fmt.Fprintf(out, "%s=%s\n", benchmarkKeyTotal, total)

	return nil
}
