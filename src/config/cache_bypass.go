package config

import cfgcache "github.com/jandedobbeleer/aliae/src/config/cache"

func SetCacheBypass(disabled bool) {
	cfgcache.SetBypass(disabled)
}

func isCacheBypassed() bool {
	return cfgcache.IsBypassed()
}
