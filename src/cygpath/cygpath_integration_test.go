package cygpath

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatchesSystemCygpathSinglePath(t *testing.T) {
	if !hasSystemCygpath() {
		t.Skip("cygpath is not available in PATH")
	}

	cases := []struct {
		name    string
		input   string
		target  Style
		command string
	}{
		{
			name:    "windows to unix",
			input:   `C:\Users\trajano\AppData\Local`,
			target:  StyleUnix,
			command: "-u",
		},
		{
			name:    "unix to windows",
			input:   "/c/Users/trajano/AppData/Local",
			target:  StyleWindows,
			command: "-w",
		},
		{
			name:    "unix to mixed",
			input:   "/c/Users/trajano/AppData/Local",
			target:  StyleMixed,
			command: "-m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.input, ConvertOptions{Target: tc.target})
			expected := mustRunSystemCygpath(t, tc.command, tc.input)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMatchesSystemCygpathPathList(t *testing.T) {
	if !hasSystemCygpath() {
		t.Skip("cygpath is not available in PATH")
	}

	cases := []struct {
		name    string
		input   string
		target  Style
		command string
	}{
		{
			name:    "path list to unix",
			input:   `C:\Users\trajano\bin;D:\tools`,
			target:  StyleUnix,
			command: "-u",
		},
		{
			name:    "path list to windows",
			input:   "/c/Users/trajano/bin:/d/tools",
			target:  StyleWindows,
			command: "-w",
		},
		{
			name:    "path list to mixed",
			input:   "/c/Users/trajano/bin:/d/tools",
			target:  StyleMixed,
			command: "-m",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Convert(tc.input, ConvertOptions{Target: tc.target, PathList: true})
			expected := mustRunSystemCygpath(t, tc.command, "-p", tc.input)
			assert.Equal(t, expected, got)
		})
	}
}

func hasSystemCygpath() bool {
	_, err := exec.LookPath("cygpath")
	return err == nil
}

func mustRunSystemCygpath(t *testing.T, args ...string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "cygpath", args...).Output()
	if err != nil {
		t.Fatalf("running cygpath failed: %v", err)
	}

	return strings.TrimSpace(string(output))
}
