package context

import "testing"

func TestIsWSL(t *testing.T) {
	originalGOOS := runtimeGOOS
	t.Cleanup(func() {
		runtimeGOOS = originalGOOS
	})

	t.Run("non-linux is false", func(t *testing.T) {
		t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
		t.Setenv("WSL_INTEROP", "/run/WSL/123_interop")
		runtimeGOOS = WINDOWS
		if isWSL() {
			t.Fatalf("expected non-linux OS to return false")
		}
	})

	t.Run("linux with wsl distro name is true", func(t *testing.T) {
		t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
		t.Setenv("WSL_INTEROP", "")
		runtimeGOOS = LINUX
		if !isWSL() {
			t.Fatalf("expected linux with WSL_DISTRO_NAME to return true")
		}
	})

	t.Run("linux with wsl interop is true", func(t *testing.T) {
		t.Setenv("WSL_DISTRO_NAME", "")
		t.Setenv("WSL_INTEROP", "/run/WSL/123_interop")
		runtimeGOOS = LINUX
		if !isWSL() {
			t.Fatalf("expected linux with WSL_INTEROP to return true")
		}
	})

	t.Run("linux without markers is false", func(t *testing.T) {
		t.Setenv("WSL_DISTRO_NAME", "")
		t.Setenv("WSL_INTEROP", "")
		runtimeGOOS = LINUX
		if isWSL() {
			t.Fatalf("expected linux without WSL env markers to return false")
		}
	})
}

func TestInitLoadsEnvironmentMap(t *testing.T) {
	t.Setenv("ALIAE_ENV_TEST", "test-value")
	Init("bash")

	if Current == nil {
		t.Fatalf("expected runtime context to be initialized")
	}

	if Current.Env["ALIAE_ENV_TEST"] != "test-value" {
		t.Fatalf("expected environment variable to be available in runtime context")
	}

	if Current.Cygpath != CygpathInternal {
		t.Fatalf("expected default cygpath mode to be internal")
	}

	if !Current.ShellLike {
		t.Fatalf("expected bash to be marked as shell-like")
	}
}

func TestIsShellLike(t *testing.T) {
	cases := []struct {
		name     string
		shell    string
		expected bool
	}{
		{name: "bash", shell: "bash", expected: true},
		{name: "zsh", shell: "zsh", expected: true},
		{name: "fish", shell: "fish", expected: true},
		{name: "tcsh", shell: "tcsh", expected: true},
		{name: "pwsh", shell: "pwsh", expected: true},
		{name: "powershell", shell: "powershell", expected: true},
		{name: "nu", shell: "nu", expected: false},
		{name: "cmd", shell: "cmd", expected: false},
		{name: "xonsh", shell: "xonsh", expected: false},
		{name: "empty", shell: "", expected: false},
	}

	for _, tc := range cases {
		if got := isShellLike(tc.shell); got != tc.expected {
			t.Fatalf("%s: expected %t, got %t", tc.name, tc.expected, got)
		}
	}
}
