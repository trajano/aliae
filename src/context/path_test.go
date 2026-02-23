package context

import (
	"os"
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
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })
	runCygpath = func(path string) (string, error) {
		if path == `C:\Users\trajano\AppData\Local\Android\Sdk\platform-tools` {
			return "/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools", nil
		}
		return "", assert.AnError
	}

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

func TestCleanPathUsesCygpath(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })

	runCygpath = func(path string) (string, error) {
		assert.Equal(t, `C:\Users\trajano\bin\`, path)
		return "/c/Users/trajano/bin", nil
	}

	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathKeepsWindowsPathWhenCygpathFails(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })

	runCygpath = func(_ string) (string, error) {
		return "", assert.AnError
	}

	assert.Equal(t, `C:\Users\trajano\bin`, cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathCachesCygpathResult(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	Current = &Runtime{OS: WINDOWS, Shell: "bash"}
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runCygpath
	t.Cleanup(func() { runCygpath = originalRunCygpath })

	calls := 0
	runCygpath = func(path string) (string, error) {
		calls++
		assert.Equal(t, `C:\Users\trajano\bin\`, path)
		return "/c/Users/trajano/bin", nil
	}

	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
	assert.Equal(t, 1, calls)
}
