package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"NoQuotes", "test", "test"},
		{"DoubleQuotes", "\"test\"", "test"},
		{"SingleQuotes", "'test'", "'test'"},
	}

	for _, tc := range tests {
		result := trimQuotes(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestReadDir(t *testing.T) {
	t.Run("ValidDir", func(t *testing.T) {
		testDirPath := filepath.Join("test", "files")
		absPath, _ := filepath.Abs(testDirPath)
		content, err := readDir(absPath)
		assert.NoError(t, err)
		assert.Equal(t, []byte("it exists\nit exists2"), content)
	})

	t.Run("NonExistentDir", func(t *testing.T) {
		_, err := readDir("path/to/nonexistent/dir")
		assert.Error(t, err)
	})
}

func TestRelativePath(t *testing.T) {
	t.Run("RelativePath", func(t *testing.T) {
		absPath, err := filepath.Abs("./test/files")
		assert.NoError(t, err)
		result, err := validatePath(absPath)
		assert.NoError(t, err)
		assert.Equal(t, absPath, result)
	})

	t.Run("Http config", func(t *testing.T) {
		configPathCache = "https://example.com/config.yaml"
		_, err := validatePath("path/to/nonex	istent/dir")
		assert.Error(t, err)
	})
}

func TestIncludeUnmarshalerCondition(t *testing.T) {
	envFile := filepath.Join(t.TempDir(), "env.yaml")
	assert.NoError(t, os.WriteFile(envFile, []byte("- name: TEST_ENV\n  value: test"), 0o600))

	aliasDir := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(aliasDir, "aliases.yaml"), []byte("- name: g\n  value: git"), 0o600))

	t.Run("true condition includes file", func(t *testing.T) {
		input := fmt.Sprintf(`env: !include %q if="eq 1 1"`, envFile)

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Contains(t, string(result), "TEST_ENV")
	})

	t.Run("false condition skips missing file without reading it", func(t *testing.T) {
		input := `env: !include "does/not/exist.yaml" if="eq 1 2"`

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("condition supports quoted template arguments", func(t *testing.T) {
		input := `env: !include "does/not/exist.yaml" if="hasCommand \"definitely-not-installed\""`

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Equal(t, "env: null", string(result))
	})

	t.Run("true condition includes directory contents", func(t *testing.T) {
		input := fmt.Sprintf(`alias: !include_dir %q if="eq 1 1"`, aliasDir)

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Contains(t, string(result), "name: g")
	})

	t.Run("false condition skips missing directory without reading it", func(t *testing.T) {
		input := `alias: !include_dir "does/not/exist" if="eq 1 2"`

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Equal(t, "alias: null", string(result))
	})

	t.Run("false condition drops an inline list item", func(t *testing.T) {
		input := "alias:\n  - !include \"does/not/exist.yaml\" if=\"eq 1 2\"\n  - name: g\n    value: git"

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Equal(t, "alias:\n\n  - name: g\n    value: git", string(result))
	})

	t.Run("path with spaces is preserved", func(t *testing.T) {
		spacedFile := filepath.Join(t.TempDir(), "my file.yaml")
		assert.NoError(t, os.WriteFile(spacedFile, []byte("- name: spaced\n  value: true"), 0o600))
		input := fmt.Sprintf(`alias: !include %q if="eq 1 1"`, spacedFile)

		result, err := includeUnmarshaler([]byte(input))

		assert.NoError(t, err)
		assert.Contains(t, string(result), "name: spaced")
	})
}

func TestIsYamlExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"YamlExtension", "test.yaml", true},
		{"YmlExtension", "test.yml", true},
		{"NoExtension", "test", false},
		{"InvalidExtension", "test.txt", false},
	}

	for _, tc := range tests {
		got := isYAMLExtension(tc.input)
		assert.Equal(t, tc.expected, got, tc.name)
	}
}
