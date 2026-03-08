package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	cfgcache "github.com/jandedobbeleer/aliae/src/config/cache"
	"github.com/jandedobbeleer/aliae/src/context"
)

func loadConfigCache(configPath string, computeVars bool) (*Aliae, bool, error) {
	cachePath := configCachePath(configPath)
	cfg, ok, err := cfgcache.Load[Aliae](cachePath, configSchemaHash(), computeVars)
	if err != nil {
		return nil, ok, err
	}
	if !ok {
		return nil, false, nil
	}
	return cfg, true, nil
}

func storeConfigCache(configPath string, computeVars bool, inputs []string, aliae *Aliae) error {
	if aliae == nil {
		return nil
	}

	cachePath := configCachePath(configPath)
	return cfgcache.Store(cachePath, configSchemaHash(), computeVars, inputs, *aliae)
}

func configCachePath(configPath string) string {
	return cfgcache.Path(configPath, cacheContextKey())
}

func cacheContextKey() string {
	current := context.GetCurrent()
	if current == nil {
		return "shell=;os=;wsl=false"
	}

	return fmt.Sprintf("shell=%s;os=%s;wsl=%t", current.Shell, current.OS, current.WSL)
}

func configSchemaHash() string {
	sum := sha256.Sum256(configSchema)
	return hex.EncodeToString(sum[:8])
}
