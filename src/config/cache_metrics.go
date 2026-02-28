package config

import "sync/atomic"

var lastLoadUsedCache atomic.Bool

func setLastLoadUsedCache(value bool) {
	lastLoadUsedCache.Store(value)
}

func LastLoadUsedCache() bool {
	return lastLoadUsedCache.Load()
}
