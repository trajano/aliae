package init

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

const renderCacheVersion = 1

var (
	envDotRefPattern     = regexp.MustCompile(`\.Env\.([A-Za-z_][A-Za-z0-9_]*)`)
	envIndexRefPattern   = regexp.MustCompile(`index\s+\.Env\s+['"]([A-Za-z_][A-Za-z0-9_]*)['"]`)
	defaultTrackedEnvKey = []string{"HOME", "USERPROFILE", "HOMEDRIVE", "HOMEPATH", "PATH"}
)

type renderCacheState struct {
	entryPath   string
	runtimeKey  string
	configPath  string
	shell       string
	trackedKeys []string
	dependency  cacheFileFingerprint
}

type cacheFileFingerprint struct {
	Path         string
	Size         int64
	ModTimeNanos int64
}

type renderCacheEntry struct {
	Script       string
	ConfigPath   string
	RuntimeKey   string
	Shell        string
	EnvHash      string
	CreatedAtUTC string
	EnvKeys      []string
	Dependency   cacheFileFingerprint
	Version      int
}

func prepareRenderCache(configPath, sh string, aliae *cfg.Aliae) *renderCacheState {
	if aliae == nil || hasStatefulScripts(aliae) {
		return nil
	}

	resolvedPath, dependencyPath, enabled, err := cfg.InitRenderCacheDependency(configPath)
	if err != nil || !enabled {
		return nil
	}

	dependency, ok := fingerprintFile(dependencyPath)
	if !ok {
		return nil
	}

	runtime := context.GetCurrent()
	if runtime == nil {
		return nil
	}

	runtimeKey := buildRuntimeKey(runtime)
	entryPath := renderCachePath(resolvedPath, sh, runtimeKey)
	trackedKeys := trackedEnvKeys(aliae)

	return &renderCacheState{
		entryPath:   entryPath,
		dependency:  dependency,
		runtimeKey:  runtimeKey,
		configPath:  resolvedPath,
		shell:       sh,
		trackedKeys: trackedKeys,
	}
}

func (s *renderCacheState) load() (string, bool) {
	if s == nil {
		return "", false
	}

	file, err := os.Open(s.entryPath)
	if err != nil {
		return "", false
	}
	defer file.Close()

	var entry renderCacheEntry
	if err := gob.NewDecoder(file).Decode(&entry); err != nil {
		return "", false
	}

	if entry.Version != renderCacheVersion || entry.ConfigPath != s.configPath || entry.RuntimeKey != s.runtimeKey || entry.Shell != s.shell {
		return "", false
	}

	if !sameFingerprint(entry.Dependency, s.dependency) {
		return "", false
	}

	if hashTrackedEnv(context.GetCurrent(), entry.EnvKeys) != entry.EnvHash {
		return "", false
	}

	return entry.Script, true
}

func (s *renderCacheState) store(script string) {
	if s == nil || strings.TrimSpace(script) == "" {
		return
	}

	entry := renderCacheEntry{
		Version:      renderCacheVersion,
		ConfigPath:   s.configPath,
		Shell:        s.shell,
		RuntimeKey:   s.runtimeKey,
		Dependency:   s.dependency,
		EnvKeys:      append([]string(nil), s.trackedKeys...),
		EnvHash:      hashTrackedEnv(context.GetCurrent(), s.trackedKeys),
		Script:       script,
		CreatedAtUTC: time.Now().UTC().Format(time.RFC3339Nano),
	}

	if err := os.MkdirAll(filepath.Dir(s.entryPath), 0o700); err != nil {
		return
	}

	file, err := os.Create(s.entryPath)
	if err != nil {
		return
	}
	defer file.Close()

	_ = gob.NewEncoder(file).Encode(entry)
}

func renderCachePath(configPath, sh, runtimeKey string) string {
	key := filepath.Clean(configPath) + "|" + sh + "|" + runtimeKey
	sum := sha256.Sum256([]byte(key))
	return aliaeState.Path(fmt.Sprintf("init-render-cache-%s.gob", hex.EncodeToString(sum[:8])))
}

func fingerprintFile(path string) (cacheFileFingerprint, bool) {
	stat, err := os.Stat(path)
	if err != nil {
		return cacheFileFingerprint{}, false
	}

	return cacheFileFingerprint{
		Path:         path,
		Size:         stat.Size(),
		ModTimeNanos: stat.ModTime().UTC().UnixNano(),
	}, true
}

func sameFingerprint(left, right cacheFileFingerprint) bool {
	return left.Path == right.Path && left.Size == right.Size && left.ModTimeNanos == right.ModTimeNanos
}

func buildRuntimeKey(runtime *context.Runtime) string {
	if runtime == nil {
		return ""
	}

	return fmt.Sprintf(
		"os=%s;wsl=%t;arch=%s;home=%s;host=%s",
		runtime.OS,
		runtime.WSL,
		runtime.Arch,
		runtime.Home,
		runtime.Hostname,
	)
}

func trackedEnvKeys(aliae *cfg.Aliae) []string {
	keys := map[string]struct{}{}
	for _, key := range defaultTrackedEnvKey {
		keys[key] = struct{}{}
	}
	if aliae == nil {
		return sortedEnvKeys(keys)
	}

	cfg.WalkConfig(aliae, cfg.ConfigVisitorFuncs{
		OnExtend: func(item *cfg.Extend) {
			collectEnvKeys(keys, string(item.If))
			collectEnvKeys(keys, item.Path)
			collectEnvKeys(keys, item.Dir)
		},
		OnVar: func(item *cfg.Var) {
			collectEnvKeys(keys, string(item.Value))
			collectEnvKeys(keys, string(item.If))
		},
		OnEnv: func(item *shell.Env) {
			collectEnvKeys(keys, item.Name)
			collectEnvKeys(keys, string(item.If))
			switch value := item.Value.(type) {
			case string:
				collectEnvKeys(keys, value)
			case []string:
				for _, line := range value {
					collectEnvKeys(keys, line)
				}
			}
		},
		OnPath: func(item *shell.Path) {
			collectEnvKeys(keys, string(item.Value))
			collectEnvKeys(keys, string(item.If))
		},
		OnCDPath: func(item *shell.CDPath) {
			collectEnvKeys(keys, string(item.Value))
			collectEnvKeys(keys, string(item.If))
		},
		OnAlias: func(item *shell.Alias) {
			collectEnvKeys(keys, item.Name)
			collectEnvKeys(keys, string(item.Value))
			collectEnvKeys(keys, string(item.If))
		},
		OnLink: func(item *shell.Link) {
			collectEnvKeys(keys, string(item.Name))
			collectEnvKeys(keys, string(item.Target))
			collectEnvKeys(keys, string(item.If))
		},
		OnScript: func(item *shell.Script) {
			collectEnvKeys(keys, string(item.Value))
			collectEnvKeys(keys, string(item.If))
			collectEnvKeys(keys, string(item.State.File))
			collectEnvKeys(keys, item.State.RunEvery)
		},
	})

	return sortedEnvKeys(keys)
}

func collectEnvKeys(set map[string]struct{}, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}

	for _, match := range envDotRefPattern.FindAllStringSubmatch(text, -1) {
		if len(match) < 2 {
			continue
		}
		set[match[1]] = struct{}{}
	}

	for _, match := range envIndexRefPattern.FindAllStringSubmatch(text, -1) {
		if len(match) < 2 {
			continue
		}
		set[match[1]] = struct{}{}
	}
}

func sortedEnvKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		if strings.TrimSpace(key) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func hashTrackedEnv(runtime *context.Runtime, keys []string) string {
	if runtime == nil || len(keys) == 0 {
		return ""
	}

	sum := sha256.New()
	for _, key := range keys {
		_, _ = sum.Write([]byte(key))
		_, _ = sum.Write([]byte{0})
		_, _ = sum.Write([]byte(runtime.Env[key]))
		_, _ = sum.Write([]byte{'\n'})
	}
	return hex.EncodeToString(sum.Sum(nil))
}

func hasStatefulScripts(aliae *cfg.Aliae) bool {
	if aliae == nil {
		return false
	}
	for _, script := range aliae.Scripts {
		if script == nil {
			continue
		}
		if strings.TrimSpace(string(script.State.File)) != "" {
			return true
		}
	}
	return false
}
