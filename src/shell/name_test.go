package shell

import "testing"

func TestShouldUseParentShell(t *testing.T) {
	t.Parallel()

	if !shouldUseParentShell("aliae", "aliae.exe") {
		t.Fatal("expected to use parent shell when process is aliae")
	}
}

func TestResolveShellNameSkipsMultipleAliaeParents(t *testing.T) {
	t.Parallel()

	parents := []string{"aliae.exe", "bash.exe"}
	i := 0
	got, _ := resolveShellName("aliae.exe", "aliae.exe", func() (string, error) {
		if i >= len(parents) {
			return "", nil
		}
		name := parents[i]
		i++
		return name, nil
	})

	if got != "bash.exe" {
		t.Fatalf("resolveShellName() = %q, want %q", got, "bash.exe")
	}
}

func TestResolveShellNameSkipsScoopShimWhenExecutableNameDiffers(t *testing.T) {
	t.Parallel()

	parents := []string{"bash.exe"}
	i := 0
	got, _ := resolveShellName("aliae-windows-amd64.exe", "aliae.exe", func() (string, error) {
		if i >= len(parents) {
			return "", nil
		}
		name := parents[i]
		i++
		return name, nil
	})

	if got != "bash.exe" {
		t.Fatalf("resolveShellName() = %q, want %q", got, "bash.exe")
	}
}
