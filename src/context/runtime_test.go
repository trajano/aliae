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
