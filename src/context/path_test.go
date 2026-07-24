package context

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDelimiterWindowsMSYS2(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash})

	assert.Equal(t, ":", PathDelimiter())
	assert.Equal(t, "/", PathSeparator())
}

func TestPathDelimiterWindowsNonMSYS2(t *testing.T) {
	t.Setenv("MSYSTEM", "")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellPwsh})

	assert.Equal(t, ";", PathDelimiter())
	assert.Equal(t, "\\", PathSeparator())
}

func TestPathContainsEquivalentWindowsAndMSYS2Forms(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	t.Setenv("PATH", "")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash, Cygpath: CygpathInternal})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	path := &Path{}
	path.Append(`C:\Users\trajano\AppData\Local\Android\Sdk\platform-tools`)

	assert.True(t, path.Contains(`/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`))

	path.Append(`/c/Users/trajano/AppData/Local/Android/Sdk/platform-tools`)
	assert.Len(t, *path, 1)
}

func TestSplitPathEntriesWindowsMSYS2MixedDelimiters(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash})

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
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellPwsh})

	original := `C:\Windows\System32;C:\Users\trajano\bin`
	t.Setenv("PATH", original)

	_ = getPath()
	assert.Equal(t, original, os.Getenv("PATH"))
}

func TestCleanPathUsesCygpath(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash, Cygpath: CygpathExternal})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runExternalCygpath
	t.Cleanup(func() { runExternalCygpath = originalRunCygpath })

	runExternalCygpath = func(path string) (string, error) {
		assert.Equal(t, `C:\Users\trajano\bin\`, path)
		return "/c/Users/trajano/bin", nil
	}

	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathKeepsWindowsPathWhenCygpathFails(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash, Cygpath: CygpathExternal})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runExternalCygpath
	t.Cleanup(func() { runExternalCygpath = originalRunCygpath })

	runExternalCygpath = func(_ string) (string, error) {
		return "", assert.AnError
	}

	assert.Equal(t, `C:\Users\trajano\bin`, cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathCachesCygpathResult(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash, Cygpath: CygpathExternal})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runExternalCygpath
	t.Cleanup(func() { runExternalCygpath = originalRunCygpath })

	calls := 0
	runExternalCygpath = func(path string) (string, error) {
		calls++
		assert.Equal(t, `C:\Users\trajano\bin\`, path)
		return "/c/Users/trajano/bin", nil
	}

	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
	assert.Equal(t, 1, calls)
}

func TestCleanPathUsesInternalCygpathByDefault(t *testing.T) {
	t.Setenv("MSYSTEM", "MINGW64")
	SetCurrent(&Runtime{OS: WINDOWS, Shell: shellBash})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunCygpath := runExternalCygpath
	t.Cleanup(func() { runExternalCygpath = originalRunCygpath })

	runExternalCygpath = func(_ string) (string, error) {
		t.Fatalf("external cygpath should not be called in internal mode")
		return "", nil
	}

	assert.Equal(t, "/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathUsesWslpathForWSL(t *testing.T) {
	SetCurrent(&Runtime{OS: LINUX, Shell: shellZsh, WSL: true})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunWslpath := runWslpath
	t.Cleanup(func() { runWslpath = originalRunWslpath })

	runWslpath = func(path string) (string, error) {
		assert.Equal(t, `C:\Users\trajano\bin\`, path)
		return "/mnt/c/Users/trajano/bin", nil
	}

	assert.Equal(t, "/mnt/c/Users/trajano/bin", cleanPath(`C:\Users\trajano\bin\`))
}

func TestCleanPathKeepsWindowsPathWhenWslpathFails(t *testing.T) {
	SetCurrent(&Runtime{OS: LINUX, Shell: shellFish, WSL: true})
	clearCleanPathCache()
	t.Cleanup(clearCleanPathCache)

	originalRunWslpath := runWslpath
	t.Cleanup(func() { runWslpath = originalRunWslpath })

	runWslpath = func(_ string) (string, error) {
		return "", assert.AnError
	}

	assert.Equal(t, `C:\Users\trajano\bin`, cleanPath(`C:\Users\trajano\bin\`))
}
