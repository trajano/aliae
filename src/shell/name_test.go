package shell

import "testing"

func TestShouldUseParentShell(t *testing.T) {
	t.Parallel()

	if !shouldUseParentShell("aliae", "aliae.exe") {
		t.Fatal("expected to use parent shell when process is aliae")
	}
}
