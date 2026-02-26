package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	text := `{{ formatString .Value}}`
	cases := []struct {
		Case     string
		Value    any
		Expected string
	}{
		{
			Case:     "string",
			Value:    "hello",
			Expected: `"hello"`,
		},
		{
			Case:     "bool",
			Value:    true,
			Expected: `true`,
		},
		{
			Case:     "int",
			Value:    32,
			Expected: `32`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: BASH}
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

// This tests both formatArray and splitString
func TestFormatArray(t *testing.T) {
	text := `{{ formatArray .Value }}`
	textDelim := `{{ formatArray .Value .Delim }}`
	cases := []struct {
		Case     string
		Value    any
		Expected string
		Delim    string
	}{
		{
			Case:     "string",
			Value:    "hello",
			Expected: `"hello"`,
		},
		{
			Case:     "Multiple Strings",
			Value:    "hello world, I am a long string",
			Expected: `"hello" "world," "I" "am" "a" "long" "string"`,
		},
		{
			Case: "Multiline String",
			Value: `hello
world
I
am
a
multiline
string`,
			Expected: `"hello" "world" "I" "am" "a" "multiline" "string"`,
		},
		{
			Case: "Single Line Starts with newline",
			Value: `
hello world I am a long string`,
			Expected: `"hello" "world" "I" "am" "a" "long" "string"`,
		},
		{
			Case:     "Single line with delimiter",
			Value:    `hello world I am a long string`,
			Delim:    ",",
			Expected: `"hello","world","I","am","a","long","string"`,
		},
		{
			Case: "Multiline with delimiter",
			Value: `hello
I
am
a
mutliline
string`,
			Delim:    ";",
			Expected: `"hello";"I";"am";"a";"mutliline";"string"`,
		},
		{
			Case:     "bool",
			Value:    true,
			Expected: `true`,
		},
		{
			Case:     "int",
			Value:    32,
			Expected: `32`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: BASH}
		var got string
		if tc.Delim == "" {
			got, _ = parse(text, tc)
		} else {
			got, _ = parse(textDelim, tc)
		}
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestEscapeString(t *testing.T) {
	text := `{{ escapeString .Value}}`
	cases := []struct {
		Case     string
		Shell    string
		Value    any
		Expected string
	}{
		{
			Case:     "string",
			Value:    `hello`,
			Expected: `hello`,
		},
		{
			Case:     "string with quotes",
			Value:    `hello "world"`,
			Expected: `hello \"world\"`,
		},
		{
			Case:     "string with backslashes",
			Value:    `hello \world`,
			Expected: `hello \\world`,
		},
		{
			Case:     "template",
			Value:    Template(`hello "world"`),
			Expected: `hello \"world\"`,
		},
		{
			Case:     "PowerShell: string",
			Shell:    PWSH,
			Value:    `hello`,
			Expected: `hello`,
		},
		{
			Case:     "PowerShell: string with quotes",
			Shell:    PWSH,
			Value:    `hello "world"`,
			Expected: "hello `\"world`\"",
		},
		{
			Case:     "PowerShell: string with backticks",
			Shell:    PWSH,
			Value:    "hello `world",
			Expected: "hello ``world",
		},
		{
			Case:     "PowerShell: template",
			Shell:    PWSH,
			Value:    Template(`hello "world"`),
			Expected: "hello `\"world`\"",
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell}
		if len(tc.Shell) == 0 {
			tc.Shell = BASH
		}
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestMatch(t *testing.T) {
	text := `{{ match .Variable "hello" "world"}}`
	cases := []struct {
		Case     string
		Variable string
		Expected string
	}{
		{
			Case:     "match",
			Variable: "hello",
			Expected: `true`,
		},
		{
			Case:     "match",
			Variable: "world",
			Expected: `true`,
		},
		{
			Case:     "noMatch",
			Variable: "goodbye",
			Expected: `false`,
		},
	}

	for _, tc := range cases {
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestTemplateHostname(t *testing.T) {
	text := `{{ .Hostname }}`
	context.Current = &context.Runtime{Shell: BASH, Hostname: "my-host"}

	got, err := parse(text, context.Current)
	assert.NoError(t, err)
	assert.Equal(t, "my-host", got)
}

func TestTemplateWSL(t *testing.T) {
	text := `{{ .WSL }}`
	context.Current = &context.Runtime{Shell: BASH, WSL: true}

	got, err := parse(text, context.Current)
	assert.NoError(t, err)
	assert.Equal(t, "true", got)
}

func TestTemplateEnv(t *testing.T) {
	text := `{{ .Env.DOTFILES }}`
	context.Current = &context.Runtime{
		Shell: BASH,
		Env: map[string]string{
			"DOTFILES": "/home/test/.dotfiles",
		},
	}

	got, err := parse(text, context.Current)
	assert.NoError(t, err)
	assert.Equal(t, "/home/test/.dotfiles", got)
}

func TestHasCommand(t *testing.T) {
	text := `{{ hasCommand .Command}}`
	t.Cleanup(clearHasCommandCache)
	clearHasCommandCache()

	cases := []struct {
		Case     string
		Command  string
		Expected string
	}{
		{
			Case:     "hasCommand",
			Command:  "go",
			Expected: `true`,
		},
		{
			Case:     "noCommand",
			Command:  "notACommand",
			Expected: `false`,
		},
	}

	for _, tc := range cases {
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestHasCommandUsesCache(t *testing.T) {
	commandName := "aliae-hascommand-cache-test"
	tempDir := t.TempDir()
	originalPath := os.Getenv("PATH")

	t.Cleanup(func() {
		_ = os.Setenv("PATH", originalPath)
		clearHasCommandCache()
	})
	clearHasCommandCache()

	missing, err := parse(`{{ hasCommand .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "false", missing)

	commandFile := filepath.Join(tempDir, commandName+executableExtension())
	content := []byte("#!/bin/sh\nexit 0\n")
	if runtime.GOOS == "windows" {
		content = []byte("@echo off\r\nexit /b 0\r\n")
	}
	assert.NoError(t, os.WriteFile(commandFile, content, 0o700))

	sep := string(os.PathListSeparator)
	assert.NoError(t, os.Setenv("PATH", tempDir+sep+originalPath))

	// Cached result remains false even though command is now available in PATH.
	stillMissing, err := parse(`{{ hasCommand .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "false", stillMissing)

	clearHasCommandCache()
	found, err := parse(`{{ hasCommand .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "true", found)
}

func TestHasCommandNoCacheBypassesCache(t *testing.T) {
	commandName := "aliae-hascommand-nocache-test"
	tempDir := t.TempDir()
	originalPath := os.Getenv("PATH")

	t.Cleanup(func() {
		_ = os.Setenv("PATH", originalPath)
		clearHasCommandCache()
	})
	clearHasCommandCache()

	missingCached, err := parse(`{{ hasCommand .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "false", missingCached)

	commandFile := filepath.Join(tempDir, commandName+executableExtension())
	content := []byte("#!/bin/sh\nexit 0\n")
	if runtime.GOOS == "windows" {
		content = []byte("@echo off\r\nexit /b 0\r\n")
	}
	assert.NoError(t, os.WriteFile(commandFile, content, 0o700))

	sep := string(os.PathListSeparator)
	assert.NoError(t, os.Setenv("PATH", tempDir+sep+originalPath))

	stillMissingCached, err := parse(`{{ hasCommand .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "false", stillMissingCached)

	foundNoCache, err := parse(`{{ hasCommandNoCache .Command }}`, struct{ Command string }{Command: commandName})
	assert.NoError(t, err)
	assert.Equal(t, "true", foundNoCache)
}

func TestFileExists(t *testing.T) {
	text := `{{ fileExists .Path }}`
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: BASH, Home: tempDir}
	relExisting := filepath.Join(tempDir, ".cache", "aliae")
	absExisting := filepath.Join(tempDir, "absolute.txt")
	assert.NoError(t, os.MkdirAll(filepath.Dir(relExisting), 0o700))
	assert.NoError(t, os.WriteFile(relExisting, []byte("ok"), 0o600))
	assert.NoError(t, os.WriteFile(absExisting, []byte("ok"), 0o600))

	t.Cleanup(clearPathExistsCache)

	cases := []struct {
		Case     string
		Path     string
		Expected string
	}{
		{
			Case:     "relative to home file exists",
			Path:     ".cache/aliae",
			Expected: "true",
		},
		{
			Case:     "absolute file exists",
			Path:     absExisting,
			Expected: "true",
		},
		{
			Case:     "relative file does not exist",
			Path:     ".cache/missing",
			Expected: "false",
		},
	}

	for _, tc := range cases {
		clearPathExistsCache()
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}

	cached := filepath.Join(tempDir, "cached.txt")
	clearPathExistsCache()
	initial, _ := parse(text, struct{ Path string }{Path: "cached.txt"})
	assert.Equal(t, "false", initial, "cached non-existent file")
	assert.NoError(t, os.WriteFile(cached, []byte("now exists"), 0o600))
	cachedResult, _ := parse(text, struct{ Path string }{Path: "cached.txt"})
	assert.Equal(t, "false", cachedResult, "should reuse cached result")
}

func TestDirExists(t *testing.T) {
	text := `{{ dirExists .Path }}`
	testDirExistsTemplate(t, text)
}

func TestFileExistsAndDirExistsShareCache(t *testing.T) {
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: BASH, Home: tempDir}
	target := filepath.Join(tempDir, "target")
	assert.NoError(t, os.MkdirAll(target, 0o700))

	t.Cleanup(clearPathExistsCache)
	clearPathExistsCache()

	fileCheck, _ := parse(`{{ fileExists .Path }}`, struct{ Path string }{Path: target})
	assert.Equal(t, "false", fileCheck)

	// dirExists should read the same cached stat result and report directory truthfully.
	dirCheck, _ := parse(`{{ dirExists .Path }}`, struct{ Path string }{Path: target})
	assert.Equal(t, "true", dirCheck)
}

func testDirExistsTemplate(t *testing.T, text string) {
	t.Helper()
	tempDir := t.TempDir()
	context.Current = &context.Runtime{Shell: BASH, Home: tempDir}
	relDir := filepath.Join(tempDir, ".cache", "aliae")
	absDir := filepath.Join(tempDir, "absolute-dir")
	assert.NoError(t, os.MkdirAll(relDir, 0o700))
	assert.NoError(t, os.MkdirAll(absDir, 0o700))

	t.Cleanup(clearPathExistsCache)

	cases := []struct {
		Case     string
		Path     string
		Expected string
	}{
		{
			Case:     "relative to home directory exists",
			Path:     ".cache/aliae",
			Expected: "true",
		},
		{
			Case:     "absolute directory exists",
			Path:     absDir,
			Expected: "true",
		},
		{
			Case:     "directory does not exist",
			Path:     ".cache/missing-dir",
			Expected: "false",
		},
	}

	for _, tc := range cases {
		clearPathExistsCache()
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}

func TestProgress(t *testing.T) {
	text := `{{ progress .Value }}`
	cases := []struct {
		Case     string
		Shell    string
		Value    any
		Expected string
	}{
		{
			Case:     "bash percentage",
			Shell:    BASH,
			Value:    71,
			Expected: `printf '\033]9;4;1;71\007'`,
		},
		{
			Case:     "bash reset",
			Shell:    BASH,
			Value:    "reset",
			Expected: `printf '\033]9;4;0;0\007'`,
		},
		{
			Case:     "pwsh percentage",
			Shell:    PWSH,
			Value:    71,
			Expected: `[Console]::Out.Write("$([char]27)]9;4;1;71$([char]7)")`,
		},
		{
			Case:     "powershell reset",
			Shell:    POWERSHELL,
			Value:    "reset",
			Expected: `[Console]::Out.Write("$([char]27)]9;4;0;0$([char]7)")`,
		},
		{
			Case:     "string percentage",
			Shell:    BASH,
			Value:    "100",
			Expected: `printf '\033]9;4;1;100\007'`,
		},
		{
			Case:     "invalid input",
			Shell:    BASH,
			Value:    "abc",
			Expected: ``,
		},
		{
			Case:     "out of range",
			Shell:    BASH,
			Value:    101,
			Expected: ``,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell}
		got, _ := parse(text, tc)
		assert.Equal(t, tc.Expected, got, tc.Case)
	}
}
