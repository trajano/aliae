package config

import (
	"strings"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

type InitBenchmarkStep struct {
	Name     string
	Duration time.Duration
}

func BenchmarkInit(configPath, sh string) ([]InitBenchmarkStep, error) {
	steps := make([]InitBenchmarkStep, 0, 12)
	record := func(name string, run func() error) error {
		start := time.Now()
		if err := run(); err != nil {
			return err
		}
		steps = append(steps, InitBenchmarkStep{
			Name:     name,
			Duration: time.Since(start),
		})
		return nil
	}

	shell.DotFile.Reset()
	defer shell.DotFile.Reset()

	var aliae *Aliae
	if err := record("init.context_init", func() error {
		context.Init(sh)
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.load_config", func() error {
		var err error
		aliae, err = LoadConfigWithoutVars(configPath)
		return err
	}); err != nil {
		return nil, err
	}

	if err := record("init.evaluate_vars", func() error {
		return aliae.ComputeVars()
	}); err != nil {
		return nil, err
	}

	if err := record("init.state_prime", func() error {
		_ = aliae.Scripts.PrimeState(time.Now())
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.autoprogress_start", func() error {
		shell.StartAutoProgress(aliae.autoProgressConfig())
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_env", func() error {
		aliae.Envs.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_path", func() error {
		aliae.Paths.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_cdpath", func() error {
		aliae.CDPaths.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_alias", func() error {
		aliae.Aliae.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_link", func() error {
		aliae.Links.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.render_script", func() error {
		aliae.Scripts.Render()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.autoprogress_end", func() error {
		shell.EndAutoProgress()
		return nil
	}); err != nil {
		return nil, err
	}

	if err := record("init.output_form", func() error {
		result := shell.DotFile.String()
		if strings.Contains(strings.ToLower(result), "aliae error:") {
			return errInitFailed
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return steps, nil
}

var errInitFailed = &initBenchmarkError{message: "init benchmark failed"}

type initBenchmarkError struct {
	message string
}

func (e *initBenchmarkError) Error() string {
	return e.message
}
