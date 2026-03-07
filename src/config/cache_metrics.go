package config

import cfgcache "github.com/jandedobbeleer/aliae/src/config/cache"

func setLastLoadUsedCache(value bool) {
	cfgcache.SetLastLoadUsed(value)
}

func LastLoadUsedCache() bool {
	return cfgcache.LastLoadUsed()
}
