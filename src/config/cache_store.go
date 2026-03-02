package config

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

const configCacheVersion = 1

type configCacheEntry struct {
	SchemaHash   string
	CreatedAtUTC string
	SourceFiles  []configCacheFile
	Config       Aliae
	Version      int
	ComputeVars  bool
}

type configCacheFile struct {
	Path         string
	Size         int64
	ModTimeNanos int64
}

func loadConfigCache(configPath string, computeVars bool) (*Aliae, bool, error) {
	cachePath := configCachePath(configPath)
	file, err := os.Open(cachePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	var entry configCacheEntry
	if err := gob.NewDecoder(file).Decode(&entry); err != nil {
		return nil, false, nil
	}

	if entry.Version != configCacheVersion || entry.SchemaHash != configSchemaHash() || entry.ComputeVars != computeVars {
		return nil, false, nil
	}

	if !isConfigCacheValid(entry.SourceFiles) {
		return nil, false, nil
	}

	cfg := entry.Config
	return &cfg, true, nil
}

func storeConfigCache(configPath string, computeVars bool, inputs []string, aliae *Aliae) error {
	if aliae == nil {
		return nil
	}

	sourceFiles, err := fingerprintConfigFiles(inputs)
	if err != nil {
		return err
	}

	cachePath := configCachePath(configPath)
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o700); err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	entry := configCacheEntry{
		Version:      configCacheVersion,
		SchemaHash:   configSchemaHash(),
		ComputeVars:  computeVars,
		Config:       *aliae,
		SourceFiles:  sourceFiles,
		CreatedAtUTC: time.Now().UTC().Format(time.RFC3339Nano),
	}

	return gob.NewEncoder(file).Encode(entry)
}

func configCachePath(configPath string) string {
	key := filepath.Clean(configPath) + "|" + cacheContextKey()
	sum := sha256.Sum256([]byte(key))
	fileName := fmt.Sprintf("config-cache-%s.gob", hex.EncodeToString(sum[:8]))
	return aliaeState.Path(fileName)
}

func cacheContextKey() string {
	if context.Current == nil {
		return "shell=;os=;wsl=false"
	}

	return fmt.Sprintf("shell=%s;os=%s;wsl=%t", context.Current.Shell, context.Current.OS, context.Current.WSL)
}

func configSchemaHash() string {
	sum := sha256.Sum256(configSchema)
	return hex.EncodeToString(sum[:8])
}

func fingerprintConfigFiles(inputs []string) ([]configCacheFile, error) {
	normalized := make([]string, 0, len(inputs))
	seen := map[string]struct{}{}
	for _, input := range inputs {
		clean := filepath.Clean(strings.TrimSpace(input))
		if clean == "" {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		normalized = append(normalized, clean)
	}
	slices.Sort(normalized)

	files := make([]configCacheFile, 0, len(normalized))
	for _, path := range normalized {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		files = append(files, configCacheFile{
			Path:         path,
			Size:         stat.Size(),
			ModTimeNanos: stat.ModTime().UTC().UnixNano(),
		})
	}

	return files, nil
}

func isConfigCacheValid(files []configCacheFile) bool {
	if len(files) == 0 {
		return false
	}

	for _, file := range files {
		stat, err := os.Stat(file.Path)
		if err != nil {
			return false
		}
		if stat.Size() != file.Size {
			return false
		}
		if stat.ModTime().UTC().UnixNano() != file.ModTimeNanos {
			return false
		}
	}

	return true
}
