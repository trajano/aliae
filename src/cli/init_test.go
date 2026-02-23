package cli

import "testing"

func TestShouldSkipInitOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ttyOnly  bool
		stdinTTY bool
		wantSkip bool
	}{
		{
			name:     "tty-only disabled",
			ttyOnly:  false,
			stdinTTY: false,
			wantSkip: false,
		},
		{
			name:     "interactive terminal",
			ttyOnly:  true,
			stdinTTY: true,
			wantSkip: false,
		},
		{
			name:     "piped output but interactive input should still print",
			ttyOnly:  true,
			stdinTTY: true,
			wantSkip: false,
		},
		{
			name:     "non-interactive process",
			ttyOnly:  true,
			stdinTTY: false,
			wantSkip: true,
		},
		{
			name:     "piped input should skip even when output is interactive",
			ttyOnly:  true,
			stdinTTY: false,
			wantSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := shouldSkipInitOutput(tt.ttyOnly, tt.stdinTTY)
			if got != tt.wantSkip {
				t.Fatalf("shouldSkipInitOutput() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}
