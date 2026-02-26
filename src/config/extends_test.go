package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigExtendsShortSyntax(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.yaml")
	child := filepath.Join(root, "child.yaml")

	require.NoError(t, os.WriteFile(base, []byte(`alias:
  - name: base
    value: from-base
`), 0o600))
	require.NoError(t, os.WriteFile(child, []byte(`extends:
  - ./base.yaml
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(child)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "base", Value: shell.Template("from-base")},
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)
}

func TestLoadConfigExtendsFailOnMissing(t *testing.T) {
	root := t.TempDir()
	missingAllowed := filepath.Join(root, "missing-allowed.yaml")
	missingBlocked := filepath.Join(root, "missing-blocked.yaml")

	require.NoError(t, os.WriteFile(missingAllowed, []byte(`extends:
  - path: ./does-not-exist.yaml
    failOnMissing: false
alias:
  - name: child
    value: from-child
`), 0o600))
	require.NoError(t, os.WriteFile(missingBlocked, []byte(`extends:
  - path: ./does-not-exist.yaml
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(missingAllowed)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)

	_, err = LoadConfig(missingBlocked)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does-not-exist.yaml")
}

func TestLoadConfigExtendsConditionalIf(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.yaml")
	child := filepath.Join(root, "child.yaml")

	require.NoError(t, os.WriteFile(base, []byte(`alias:
  - name: base
    value: from-base
`), 0o600))
	require.NoError(t, os.WriteFile(child, []byte(`extends:
  - path: ./base.yaml
    if: 'false'
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(child)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)
}

func TestLoadConfigExtendsConditionalIfTrue(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.yaml")
	child := filepath.Join(root, "child.yaml")

	require.NoError(t, os.WriteFile(base, []byte(`alias:
  - name: base
    value: from-base
`), 0o600))
	require.NoError(t, os.WriteFile(child, []byte(`extends:
  - path: ./base.yaml
    if: 'true'
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(child)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "base", Value: shell.Template("from-base")},
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)
}

func TestLoadConfigExtendsDirectory(t *testing.T) {
	root := t.TempDir()
	parts := filepath.Join(root, "parts")
	require.NoError(t, os.MkdirAll(filepath.Join(parts, "sub"), 0o700))

	require.NoError(t, os.WriteFile(filepath.Join(parts, "A.yaml"), []byte(`alias:
  - name: A
    value: from-A
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(parts, "a.yml"), []byte(`alias:
  - name: a
    value: from-a
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(parts, "sub", "z.yaml"), []byte(`alias:
  - name: z
    value: from-z
`), 0o600))

	nonRecursive := filepath.Join(root, "non-recursive.yaml")
	require.NoError(t, os.WriteFile(nonRecursive, []byte(`extends:
  - dir: ./parts
    recursive: false
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(nonRecursive)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "A", Value: shell.Template("from-A")},
		{Name: "a", Value: shell.Template("from-a")},
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)

	recursive := filepath.Join(root, "recursive.yaml")
	require.NoError(t, os.WriteFile(recursive, []byte(`extends:
  - dir: ./parts
    recursive: true
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err = LoadConfig(recursive)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "A", Value: shell.Template("from-A")},
		{Name: "a", Value: shell.Template("from-a")},
		{Name: "z", Value: shell.Template("from-z")},
		{Name: "child", Value: shell.Template("from-child")},
	}, got.Aliae)
}

func TestLoadConfigExtendsDirectoryWithNestedExtends(t *testing.T) {
	root := t.TempDir()
	parts := filepath.Join(root, "parts")
	fragments := filepath.Join(root, "fragments")
	require.NoError(t, os.MkdirAll(parts, 0o700))
	require.NoError(t, os.MkdirAll(fragments, 0o700))

	require.NoError(t, os.WriteFile(filepath.Join(fragments, "shared.yaml"), []byte(`alias:
  - name: shared
    value: from-shared
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(parts, "10-base.yaml"), []byte(`alias:
  - name: base
    value: from-base
`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(parts, "20-child.yaml"), []byte(`extends:
  - ../fragments/shared.yaml
alias:
  - name: child
    value: from-child
`), 0o600))

	config := filepath.Join(root, "root.yaml")
	require.NoError(t, os.WriteFile(config, []byte(`extends:
  - dir: ./parts
alias:
  - name: root
    value: from-root
`), 0o600))

	got, err := LoadConfig(config)
	require.NoError(t, err)
	assert.Equal(t, shell.Aliae{
		{Name: "base", Value: shell.Template("from-base")},
		{Name: "shared", Value: shell.Template("from-shared")},
		{Name: "child", Value: shell.Template("from-child")},
		{Name: "root", Value: shell.Template("from-root")},
	}, got.Aliae)
}

func TestLoadConfigExtendsCycle(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "a.yaml")
	b := filepath.Join(root, "b.yaml")

	require.NoError(t, os.WriteFile(a, []byte(`extends:
  - ./b.yaml
alias:
  - name: a
    value: a
`), 0o600))
	require.NoError(t, os.WriteFile(b, []byte(`extends:
  - ./a.yaml
alias:
  - name: b
    value: b
`), 0o600))

	_, err := LoadConfig(a)
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "cycle")
}

func TestLoadConfigExtendsDepthLimit(t *testing.T) {
	root := t.TempDir()
	const files = 11

	for i := 1; i <= files; i++ {
		file := filepath.Join(root, fmt.Sprintf("c%d.yaml", i))
		next := ""
		if i < files {
			next = fmt.Sprintf("c%d.yaml", i+1)
		}

		content := ""
		if len(next) > 0 {
			content += "extends:\n  - ./" + next + "\n"
		}

		content += fmt.Sprintf("alias:\n  - name: file%d\n    value: ok\n", i)
		require.NoError(t, os.WriteFile(file, []byte(content), 0o600))
	}

	_, err := LoadConfig(filepath.Join(root, "c1.yaml"))
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "depth")
}

func TestLoadConfigExtendsIgnoresProgressInternalFromIncludedFiles(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.yaml")
	child := filepath.Join(root, "child.yaml")

	require.NoError(t, os.WriteFile(base, []byte(`progress:
  start_percentage: 0
  internal: 33
alias:
  - name: base
    value: from-base
`), 0o600))

	require.NoError(t, os.WriteFile(child, []byte(`extends:
  - ./base.yaml
progress:
  start_percentage: 10
  internal: 7
alias:
  - name: child
    value: from-child
`), 0o600))

	got, err := LoadConfig(child)
	require.NoError(t, err)
	assert.True(t, got.Progress.Enabled)
	assert.Equal(t, 10.0, got.Progress.StartPercentage)
	assert.Equal(t, 7.0, got.Progress.Internal)
}
