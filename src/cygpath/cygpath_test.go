package cygpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowsToUnix(t *testing.T) {
	path, ok := WindowsToUnix(`C:\Users\trajano\AppData\Local`)
	assert.True(t, ok)
	assert.Equal(t, "/c/Users/trajano/AppData/Local", path)
}

func TestWindowsToUnixUNC(t *testing.T) {
	path, ok := WindowsToUnix(`\\server\share\tools`)
	assert.True(t, ok)
	assert.Equal(t, "//server/share/tools", path)
}

func TestUnixToWindows(t *testing.T) {
	path, ok := UnixToWindows("/c/Users/trajano/AppData/Local")
	assert.True(t, ok)
	assert.Equal(t, `C:\Users\trajano\AppData\Local`, path)
}

func TestUnixToWindowsUNC(t *testing.T) {
	path, ok := UnixToWindows("//server/share/tools")
	assert.True(t, ok)
	assert.Equal(t, `\\server\share\tools`, path)
}

func TestConvertMixed(t *testing.T) {
	assert.Equal(t, "C:/Users/trajano/tools", Convert("/c/Users/trajano/tools", ConvertOptions{Target: StyleMixed}))
}

func TestConvertListToUnix(t *testing.T) {
	input := `C:\Users\trajano\bin;D:\tools:/c/Users/trajano/.local/bin`
	got := Convert(input, ConvertOptions{Target: StyleUnix, PathList: true})

	assert.Equal(t, "/c/Users/trajano/bin:/d/tools:/c/Users/trajano/.local/bin", got)
}

func TestSplitList(t *testing.T) {
	paths := SplitList(`C:\Users\trajano\bin;C:\Program Files\Git\usr\bin:/c/Users/trajano/.local/bin`)
	assert.Equal(t, []string{
		`C:\Users\trajano\bin`,
		`C:\Program Files\Git\usr\bin`,
		`/c/Users/trajano/.local/bin`,
	}, paths)
}
