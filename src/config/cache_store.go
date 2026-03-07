package config

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	cfgcache "github.com/jandedobbeleer/aliae/src/config/cache"
	"github.com/jandedobbeleer/aliae/src/context"
)

func loadConfigCache(configPath string, computeVars bool) (*Aliae, bool, error) {
	cachePath := configCachePath(configPath)
	payload, ok, err := cfgcache.Load(cachePath, configSchemaHash(), computeVars)
	if err != nil {
		return nil, ok, err
	}
	if !ok {
		return nil, false, nil
	}

	var cfg Aliae
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&cfg); err != nil {
		return nil, false, nil
	}

	return &cfg, true, nil
}

func storeConfigCache(configPath string, computeVars bool, inputs []string, aliae *Aliae) error {
	if aliae == nil {
		return nil
	}

	var payload bytes.Buffer
	if err := gob.NewEncoder(&payload).Encode(*aliae); err != nil {
		return err
	}

	cachePath := configCachePath(configPath)
	return cfgcache.Store(cachePath, configSchemaHash(), computeVars, inputs, payload.Bytes())
}

func configCachePath(configPath string) string {
	return cfgcache.Path(configPath, cacheContextKey())
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
