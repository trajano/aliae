package config

import "sync/atomic"

var cacheBypass atomic.Bool

func SetCacheBypass(disabled bool) {
	cacheBypass.Store(disabled)
}

func isCacheBypassed() bool {
	return cacheBypass.Load()
}
