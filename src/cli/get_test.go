package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetShellWritesToStdout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := getCmd.RunE(cmd, []string{"shell"})
	require.NoError(t, err)
	assert.NotEmpty(t, strings.TrimSpace(stdout.String()))
	assert.Empty(t, stderr.String())
}
