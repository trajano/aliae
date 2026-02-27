package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
)

func TestResolveConfigPath(t *testing.T) {
	cases := []struct {
		name      string
		configVar string
		homeVar   string
		expected  string
	}{
		{"Config env var", "test", "", "test"},
		{"No config env var", "", "/home", "/home/.aliae.yaml"},
		{"No config env var, no home", "", "", ".aliae.yaml"},
	}

	for _, c := range cases {
		os.Setenv("ALIAE_CONFIG", c.configVar)
		os.Setenv("HOME", c.homeVar)
		got := resolveConfigPath("")
		assert.Equal(t, got, got, c.name)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		expected    *Aliae
		name        string
		config      string
		expectError bool
	}{
		{
			name:   "Valid",
			config: "aliae.valid.yaml",
			expected: &Aliae{
				Aliae: shell.Aliae{
					{Name: "test", Value: shell.Template("test")},
					{Name: "test2", Value: shell.Template("test2")},
				},
				Envs: shell.Envs{
					{Name: "TEST_ENV", Value: "test"},
				},
			},
		},
		{
			name:        "Invalid",
			config:      "aliae.invalid.yaml",
			expectError: true,
		},
		{
			name:   "Valid with generic",
			config: "aliae.valid_generic.yaml",
			expected: &Aliae{
				Aliae: shell.Aliae{
					{Name: "test", Value: shell.Template("test")},
					{Name: "test2", Value: shell.Template("test2")},
					{Name: "test3", Value: shell.Template("test3")},
				},
			},
		},
	}

	for _, tc := range tests {
		configFile := filepath.Join("test", tc.config)
		got, err := LoadConfig(configFile)

		if tc.expectError {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}

		assert.Equal(t, tc.expected, got, tc.name)
	}
}

func TestResolveConfigDir(t *testing.T) {
	cases := []struct {
		name     string
		config   string
		expected string
	}{
		{
			name:     "Local path",
			config:   "/tmp/aliae/aliae.yaml",
			expected: filepath.Dir("/tmp/aliae/aliae.yaml"),
		},
		{
			name:     "Remote URL",
			config:   "https://example.com/configs/aliae.yaml?x=1",
			expected: "https://example.com/configs",
		},
	}

	for _, tc := range cases {
		got := resolveConfigDir(tc.config)
		assert.Equal(t, tc.expected, got, tc.name)
	}
}

func TestParseConfigStatTimeout(t *testing.T) {
	aliae, err := parseConfig([]byte("stat_timeout: 250ms\n"), true)
	assert.NoError(t, err)
	assert.Equal(t, 250*time.Millisecond, aliae.StatTimeout)
}

func TestResolveTemplateContext(t *testing.T) {
	t.Setenv("ALIAE_CONFIG", "")
	t.Setenv("HOME", "/home/user")

	path, dir := ResolveTemplateContext("")
	assert.Equal(t, "/home/user/.aliae.yaml", path)
	assert.Equal(t, filepath.Dir("/home/user/.aliae.yaml"), dir)
}

func TestParseConfigRejectsInvalidScriptStateRunEvery(t *testing.T) {
	_, err := parseConfig([]byte(`script:
  - value: echo hello
    state:
      file: hello.state
      runEvery: nope
`), true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "script[0].state.runEvery")
}

func TestParseConfigRejectsDuplicateScriptStateFile(t *testing.T) {
	_, err := parseConfig([]byte(`script:
  - value: echo hello
    state:
      file: shared.state
  - value: echo bye
    state:
      file: shared.state
`), true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicates")
}

func TestParseConfigRejectsVarUsingVarTemplateValue(t *testing.T) {
	_, err := parseConfig([]byte(`var:
  - name: A
    value: '{{ .Var.B }}'
`), true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "var[0].value")
}
