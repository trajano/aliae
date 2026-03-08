package config

func InitRenderCacheDependency(configPath string) (resolvedPath, cachePath string, enabled bool, err error) {
	resolvedPath = resolveConfigPath(configPath)
	enabled, err = loadRootCache(resolvedPath)
	if err != nil {
		return resolvedPath, "", false, err
	}

	if isCacheBypassed() {
		enabled = false
	}
	if !enabled {
		return resolvedPath, "", false, nil
	}

	return resolvedPath, configCachePath(resolvedPath), true, nil
}
