package context

import "os"

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	DARWIN  = "darwin"
)

func isMSYS2Shell() bool {
	current := GetCurrent()
	if current == nil {
		return false
	}

	return current.OS == WINDOWS && current.Shell == "bash" && os.Getenv("MSYSTEM") != ""
}

func PathDelimiter() string {
	current := GetCurrent()
	if current == nil {
		return ":"
	}

	switch current.OS {
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
	current := GetCurrent()
	if current == nil {
		return "/"
	}

	switch current.OS {
	case WINDOWS:
		if isMSYS2Shell() {
			return "/"
		}
		return "\\"
	default:
		return "/"
	}
}
