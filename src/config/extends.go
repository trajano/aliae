package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	contextpkg "github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

const maxExtendsDepth = 10

type Extend struct {
	FailOnMissing *bool    `yaml:"failOnMissing"`
	Path          string   `yaml:"path"`
	Dir           string   `yaml:"dir"`
	If            shell.If `yaml:"if"`
	Recursive     bool     `yaml:"recursive"`
}

func (e *Extend) UnmarshalYAML(data []byte) error {
	var path string
	if err := yaml.Unmarshal(data, &path); err == nil {
		e.Path = path
		return nil
	}

	type rawExtend Extend
	var raw rawExtend
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	*e = Extend(raw)
	return nil
}

func (e Extend) shouldFailOnMissing() bool {
	return e.FailOnMissing == nil || *e.FailOnMissing
}

func loadLocalConfigWithInputs(configPath string) (*Aliae, []string, error) {
	return loadLocalConfigRecursive(configPath, nil, 1)
}

func loadLocalConfigRecursive(configPath string, stack []string, depth int) (*Aliae, []string, error) {
	if depth > maxExtendsDepth {
		return nil, nil, fmt.Errorf("extends depth limit exceeded: max %d", maxExtendsDepth)
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		absPath = configPath
	}

	absPath = filepath.Clean(absPath)
	if slices.Contains(stack, absPath) {
		return nil, nil, fmt.Errorf("extends cycle detected for %s", absPath)
	}

	if stat, statErr := os.Stat(absPath); os.IsNotExist(statErr) || (statErr == nil && stat.IsDir()) {
		return nil, nil, fmt.Errorf("config file not found: %s", absPath)
	}

	previousConfigPath := configPathCache
	configPathCache = absPath
	defer func() {
		configPathCache = previousConfigPath
	}()

	previousTemplatePath := ""
	previousTemplateDir := ""
	runtime := contextpkg.GetCurrent()
	if runtime != nil {
		previousTemplatePath = runtime.ConfigPath
		previousTemplateDir = runtime.ConfigDir
	}
	setTemplateConfigContext(absPath)
	defer func() {
		runtime := contextpkg.GetCurrent()
		if runtime == nil {
			return
		}

		runtime.ConfigPath = previousTemplatePath
		runtime.ConfigDir = previousTemplateDir
	}()

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, nil, err
	}

	includedData, err := includeUnmarshaler(data)
	if err != nil {
		return nil, nil, err
	}
	if err := validateScriptWeightsInYAML(includedData); err != nil {
		return nil, nil, err
	}
	if err := validateScriptStateInYAML(includedData); err != nil {
		return nil, nil, err
	}

	extends, err := parseExtends(includedData)
	if err != nil {
		return nil, nil, err
	}
	if depth == 1 {
		MarkInitInternalProgressDiscoveryComplete()
	}

	merged := &Aliae{}
	inputs := []string{absPath}
	nextStack := make([]string, len(stack)+1)
	copy(nextStack, stack)
	nextStack[len(stack)] = absPath
	for _, item := range extends {
		paths, pathErr := resolveExtendsPaths(item, absPath)
		if pathErr != nil {
			return nil, nil, pathErr
		}

		for _, path := range paths {
			parent, parentInputs, loadErr := loadLocalConfigRecursive(path, nextStack, depth+1)
			if loadErr != nil {
				return nil, nil, loadErr
			}
			parent = applyExtendsCondition(parent, item.If)

			merged.merge(parent)
			inputs = append(inputs, parentInputs...)
			if depth == 1 {
				MarkInitInternalProgressLinkedConfigLoaded()
			}
		}
	}

	currentConfig, err := decodeAliae(includedData)
	if err != nil {
		return nil, nil, err
	}

	merged.merge(currentConfig)

	return merged, inputs, nil
}

func decodeAliae(data []byte) (*Aliae, error) {
	var aliae Aliae

	decoder := yaml.NewDecoder(
		bytes.NewBuffer(data),
		yaml.CustomUnmarshaler(templateUmarshaler),
		yaml.CustomUnmarshaler(progressUnmarshaler),
	)
	if err := decoder.Decode(&aliae); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	if err := hydrateProgressFromYAML(&aliae, data); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return &aliae, nil
}

func parseExtends(data []byte) ([]Extend, error) {
	var doc struct {
		Extends []Extend `yaml:"extends"`
	}

	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse extends block: %s", err)
	}

	return doc.Extends, nil
}

func resolveExtendsPaths(item Extend, configPath string) ([]string, error) {
	hasPath := len(strings.TrimSpace(item.Path)) > 0
	hasDir := len(strings.TrimSpace(item.Dir)) > 0

	if hasPath == hasDir {
		return nil, errors.New("extends entry must set exactly one of path or dir")
	}

	if hasPath {
		path, err := resolveRelativePath(item.Path, configPath)
		if err != nil {
			if !item.shouldFailOnMissing() && errors.Is(err, fs.ErrNotExist) {
				return nil, nil
			}

			return nil, err
		}

		stat, err := os.Stat(path)
		if err != nil {
			if !item.shouldFailOnMissing() && errors.Is(err, fs.ErrNotExist) {
				return nil, nil
			}

			return nil, err
		}

		if stat.IsDir() {
			return nil, fmt.Errorf("extends path is a directory: %s", path)
		}

		return []string{path}, nil
	}

	dir, err := resolveRelativePath(item.Dir, configPath)
	if err != nil {
		if !item.shouldFailOnMissing() && errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	files, err := readExtendsDir(dir, item.Recursive)
	if err != nil {
		if !item.shouldFailOnMissing() && errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	return files, nil
}

func readExtendsDir(dir string, recursive bool) ([]string, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("extends dir is not a directory: %s", dir)
	}

	files := make([]string, 0)
	if recursive {
		walkErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !isYAMLExtension(d.Name()) {
				return nil
			}

			files = append(files, path)
			return nil
		})
		if walkErr != nil {
			return nil, walkErr
		}
	} else {
		entries, readErr := os.ReadDir(dir)
		if readErr != nil {
			return nil, readErr
		}

		for _, entry := range entries {
			if entry.IsDir() || !isYAMLExtension(entry.Name()) {
				continue
			}

			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	slices.Sort(files)
	return files, nil
}

func (a *Aliae) merge(other *Aliae) {
	if other == nil {
		return
	}

	a.Aliae = append(a.Aliae, other.Aliae...)
	a.Extends = append(a.Extends, other.Extends...)
	a.Vars = append(a.Vars, other.Vars...)
	a.Envs = append(a.Envs, other.Envs...)
	a.Paths = append(a.Paths, other.Paths...)
	a.CDPaths = append(a.CDPaths, other.CDPaths...)
	a.Scripts = append(a.Scripts, other.Scripts...)
	a.Links = append(a.Links, other.Links...)

	if other.StatTimeout > 0 {
		a.StatTimeout = other.StatTimeout
	}

	if len(strings.TrimSpace(other.Cygpath)) > 0 {
		a.Cygpath = other.Cygpath
	}

	if other.Progress.Enabled || other.Progress.EndPercentage.Reset || other.Progress.EndPercentage.Value > 0 || other.Progress.StartPercentage > 0 {
		a.Progress = other.Progress
	}
}

func applyExtendsCondition(aliae *Aliae, condition shell.If) *Aliae {
	if aliae == nil {
		return nil
	}

	if strings.TrimSpace(string(condition)) == "" {
		return aliae
	}

	WalkConfig(aliae, ConfigVisitorFuncs{
		OnAlias: func(alias *shell.Alias) {
			alias.If = andIf(alias.If, condition)
		},
		OnVar: func(variable *Var) {
			variable.If = andIf(variable.If, condition)
		},
		OnEnv: func(env *shell.Env) {
			env.If = andIf(env.If, condition)
		},
		OnPath: func(path *shell.Path) {
			path.If = andIf(path.If, condition)
		},
		OnCDPath: func(cdpath *shell.CDPath) {
			cdpath.If = andIf(cdpath.If, condition)
		},
		OnScript: func(script *shell.Script) {
			script.If = andIf(script.If, condition)
		},
		OnLink: func(link *shell.Link) {
			link.If = andIf(link.If, condition)
		},
	})

	return aliae
}

func andIf(existing, condition shell.If) shell.If {
	if strings.TrimSpace(string(existing)) == "" {
		return condition
	}

	if strings.TrimSpace(string(condition)) == "" {
		return existing
	}

	return shell.If(fmt.Sprintf("and (%s) (%s)", condition, existing))
}
