package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDelimiterWindowsMSYS2(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}

	assert.Equal(t, ":", PathDelimiter())
	assert.Equal(t, "/", PathSeparator())
}

func TestPathDelimiterWindowsNonMSYS2(t *testing.T) {
	t.Setenv("MSYSTEM", "")
	Current = &Runtime{OS: WINDOWS, Shell: "pwsh"}

	assert.Equal(t, ";", PathDelimiter())
	assert.Equal(t, "\\", PathSeparator())
}

func TestPathContainsEquivalentWindowsAndMSYS2Forms(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	t.Setenv("PATH", "")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}

	path := &Path{}
	path.Append(`C:\Users\trajano\AppData\Local\Android\Sdk\platform-tools`)

	assert.True(t, path.Contains(`/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`))

	path.Append(`/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`)
	assert.Len(t, *path, 1)
}
