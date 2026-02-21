package context

import "os"

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	DARWIN  = "darwin"
)

func isMSYS2Shell() bool {
	if Current == nil {
		return false
	}

	return Current.OS == WINDOWS && Current.Shell == "bash" && os.Getenv("MSYSTEM") != ""
}

func PathDelimiter() string {
	if Current == nil {
		return ":"
	}

	switch Current.OS {
	case WINDOWS:
		if isMSYS2Shell() {
			return ":"
		}
		return ";"
	default:
		return ":"
	}
}

func PathSeparator() string {
	if Current == nil {
		return "/"
	}

	switch Current.OS {
	case WINDOWS:
		if isMSYS2Shell() {
			return "/"
		}
		return "\\"
	default:
		return "/"
	}
}
