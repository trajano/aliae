package load

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveTemplateContext(t *testing.T) {
	t.Setenv("ALIAE_CONFIG", "")
	t.Setenv("HOME", "/home/user")

	path, dir := ResolveTemplateContext("")
	assert.Equal(t, "/home/user/.aliae.yaml", path)
	assert.Equal(t, filepath.Dir("/home/user/.aliae.yaml"), dir)
}

func TestResolveConfigDirRemote(t *testing.T) {
	got := ResolveConfigDir("https://example.com/configs/aliae.yaml?x=1")
	assert.Equal(t, "https://example.com/configs", got)
}
