package cache

import "sync/atomic"

var lastLoadUsed atomic.Bool
var bypass atomic.Bool

func SetLastLoadUsed(value bool) {
	lastLoadUsed.Store(value)
}

func LastLoadUsed() bool {
	return lastLoadUsed.Load()
}

func SetBypass(disabled bool) {
	bypass.Store(disabled)
}

func IsBypassed() bool {
	return bypass.Load()
}
