package cygpath

import (
	"strings"
)

type Style string

const (
	StyleUnix    Style = "unix"
	StyleWindows Style = "windows"
	StyleMixed   Style = "mixed"
)

type ConvertOptions struct {
	Target   Style
	PathList bool
}

func Convert(input string, options ConvertOptions) string {
	if options.PathList {
		return ConvertList(input, options.Target)
	}

	switch options.Target {
	case StyleUnix:
		return ToUnix(input)
	case StyleWindows:
		return ToWindows(input)
	case StyleMixed:
		return ToMixed(input)
	default:
		return input
	}
}

func ConvertList(input string, target Style) string {
	paths := SplitList(input)
	converted := make([]string, len(paths))
	for i, path := range paths {
		converted[i] = Convert(path, ConvertOptions{Target: target})
	}

	return strings.Join(converted, listDelimiter(target))
}

func SplitList(input string) []string {
	if input == "" {
		return []string{}
	}

	result := make([]string, 0)
	start := 0

	for i := range len(input) {
		switch input[i] {
		case ';':
			result = append(result, input[start:i])
			start = i + 1
		case ':':
			if i == start+1 && start < len(input) && isASCIIAlpha(input[start]) {
				continue
			}
			result = append(result, input[start:i])
			start = i + 1
		}
	}

	result = append(result, input[start:])
	return result
}

func ToUnix(path string) string {
	if converted, ok := WindowsToUnix(path); ok {
		return converted
	}

	return strings.ReplaceAll(path, `\`, "/")
}

func ToWindows(path string) string {
	if converted, ok := UnixToWindows(path); ok {
		return converted
	}

	return strings.ReplaceAll(path, "/", `\`)
}

func ToMixed(path string) string {
	return strings.ReplaceAll(ToWindows(path), `\`, "/")
}

func WindowsToUnix(path string) (string, bool) {
	if path == "" {
		return "", false
	}

	if strings.HasPrefix(path, `\\`) {
		return "//" + strings.ReplaceAll(strings.TrimPrefix(path, `\\`), `\`, "/"), true
	}

	if len(path) < 3 || path[1] != ':' {
		return "", false
	}

	drive := path[0]
	if !isASCIIAlpha(drive) {
		return "", false
	}

	if path[2] != '\\' && path[2] != '/' {
		return "", false
	}

	rest := strings.ReplaceAll(path[2:], `\`, `/`)
	return "/" + strings.ToLower(string(drive)) + rest, true
}

func UnixToWindows(path string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}

	if strings.HasPrefix(path, "//") {
		rest := strings.TrimPrefix(path, "//")
		if rest == "" {
			return "", false
		}
		return `\\` + strings.ReplaceAll(rest, `/`, `\`), true
	}

	if len(path) < 4 || path[0] != '/' || path[2] != '/' {
		return "", false
	}

	drive := path[1]
	if !isASCIIAlpha(drive) {
		return "", false
	}

	rest := strings.ReplaceAll(path[3:], `/`, `\`)
	if rest == "" {
		return strings.ToUpper(string(drive)) + `:\`, true
	}

	return strings.ToUpper(string(drive)) + `:\` + rest, true
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
}

func listDelimiter(target Style) string {
	switch target {
	case StyleUnix:
		return ":"
	case StyleWindows, StyleMixed:
		return ";"
	default:
		return ":"
	}
}
