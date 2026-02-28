package config

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

func applyCygpathDefaults(aliae *Aliae) {
	if aliae == nil {
		return
	}

	aliae.Cygpath = context.NormalizeCygpathMode(aliae.Cygpath)
}

func validateCygpathMode(aliae *Aliae) error {
	if aliae == nil {
		return nil
	}

	switch aliae.Cygpath {
	case context.CygpathInternal, context.CygpathExternal:
		return nil
	default:
		return fmt.Errorf("invalid cygpath: %q (must be internal or external)", aliae.Cygpath)
	}
}

func setRuntimeCygpathMode(aliae *Aliae) {
	if aliae == nil || context.Current == nil {
		return
	}

	context.Current.Cygpath = aliae.Cygpath
}
