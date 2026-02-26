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
}
