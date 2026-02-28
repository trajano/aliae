package config

import (
	stdcontext "context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-yaml"
)

func loadRootProgress(configPath string) (Progress, error) {
	document, err := loadRootConfigDocument(configPath)
	if err != nil {
		return Progress{}, err
	}

	progressValue, hasProgress := document["progress"]
	if !hasProgress {
		return Progress{}, nil
	}

	return parseProgressAny(progressValue)
}

func loadRootCache(configPath string) (bool, error) {
	document, err := loadRootConfigDocument(configPath)
	if err != nil {
		return false, err
	}

	cacheValue, hasCache := document["cache"]
	if !hasCache {
		return true, nil
	}

	cache, ok := cacheValue.(bool)
	if !ok {
		return false, fmt.Errorf("cache must be a boolean")
	}

	return cache, nil
}

func loadRootConfigDocument(configPath string) (map[string]any, error) {
	data, err := readRootConfigBytes(configPath)
	if err != nil {
		return nil, err
	}

	document := map[string]any{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	return document, nil
}

func readRootConfigBytes(configPath string) ([]byte, error) {
	if isRemoteConfigPath(configPath) {
		req, err := http.NewRequestWithContext(stdcontext.Background(), "GET", configPath, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download config file: %s", resp.Status)
		}

		return io.ReadAll(resp.Body)
	}

	return os.ReadFile(configPath)
}
