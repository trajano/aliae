package context

import (
	"os"
	"os/exec"
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

func TestSplitPathEntriesWindowsMSYS2MixedDelimiters(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}

	paths := `C:\Users\trajano\bin;C:\Program Files\Git\usr\bin:/c/Users/trajano/.local/bin:/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`
	got := splitPathEntries(paths)

	assert.Equal(t, []string{
		`C:\Users\trajano\bin`,
		`C:\Program Files\Git\usr\bin`,
		`/c/Users/trajano/.local/bin`,
		`/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`,
	}, got)
}

func TestGetPathDoesNotMutateEnvironmentPath(t *testing.T) {
	t.Setenv("MSYSTEM", "")
	Current = &Runtime{OS: WINDOWS, Shell: "pwsh"}

	original := `C:\Windows\System32;C:\Users\trajano\bin`
	t.Setenv("PATH", original)

	_ = getPath()
	assert.Equal(t, original, os.Getenv("PATH"))
}

func TestWindowsToMSYSPathUsesCygpath(t *testing.T) {
	if _, err := exec.LookPath("cygpath"); err != nil {
		t.Skip("cygpath not available")
	}

	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })

	runCygpath = func(path string) (string, error) {
		assert.Equal(t, `C:\Users\trajano\bin`, path)
		return "/c/Users/trajano/bin", nil
	}

	got, ok := windowsToMSYSPath(`C:\Users\trajano\bin`)
	assert.True(t, ok)
	assert.Equal(t, "/c/Users/trajano/bin", got)
}

func TestWindowsToMSYSPathFallsBackWhenCygpathFails(t *testing.T) {
	if _, err := exec.LookPath("cygpath"); err != nil {
		t.Skip("cygpath not available")
	}

	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })

	runCygpath = func(_ string) (string, error) {
		return "", assert.AnError
	}

	got, ok := windowsToMSYSPath(`C:\Users\trajano\bin`)
	assert.True(t, ok)
	assert.Equal(t, "/c/Users/trajano/bin", got)
}
