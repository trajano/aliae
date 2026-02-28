package context

import "strings"

const (
	CygpathInternal = "internal"
	CygpathExternal = "external"
)

func NormalizeCygpathMode(mode string) string {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	if normalized == "" {
		return CygpathInternal
	}

	switch normalized {
	case CygpathInternal, CygpathExternal:
		return normalized
	default:
		return normalized
	}
}
