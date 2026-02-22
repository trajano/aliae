package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestCDPath(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		CDPath   *CDPath
		OS       string
		Expected string
	}{
		{
			Case:  "Unknown shell",
			Shell: "FOO",
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
		},
		{
			Case:  "PWSH is skipped",
			Shell: PWSH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
		},
		{
			Case:  "CMD is skipped",
			Shell: CMD,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
		},
		{
			Case:  "NU is skipped",
			Shell: NU,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
		},
		{
			Case:  "BASH - single item",
			Shell: BASH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: `export CDPATH="/usr/local/share:$CDPATH"`,
		},
		{
			Case:  "BASH - multiple items",
			Shell: BASH,
			CDPath: &CDPath{
				Value: "/usr/local/share\n/usr/share",
			},
			Expected: `export CDPATH="/usr/local/share:$CDPATH"
export CDPATH="/usr/share:$CDPATH"`,
		},
		{
			Case:  "ZSH - single item",
			Shell: ZSH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: `cdpath=( /usr/local/share $cdpath )`,
		},
		{
			Case:  "FISH - single item",
			Shell: FISH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: `set -g cdpath /usr/local/share $cdpath`,
		},
		{
			Case:  "TCSH - single item",
			Shell: TCSH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: `set cdpath = ( /usr/local/share $cdpath );`,
		},
		{
			Case:  "XONSH - single item",
			Shell: XONSH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: `$CDPATH = ["/usr/local/share"] + $CDPATH`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{
			Shell:  tc.Shell,
			Home:   "/Users/jan",
			OS:     tc.OS,
			CDPath: &context.Path{"/opt/cd"},
		}
		assert.Equal(t, tc.Expected, tc.CDPath.string(), tc.Case)
	}
}

func TestCDPathForce(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		CDPath   *CDPath
		Expected string
	}{
		{
			Case:  "BASH - Not Force",
			Shell: BASH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
			},
			Expected: "",
		},
		{
			Case:  "BASH - Force",
			Shell: BASH,
			CDPath: &CDPath{
				Value: "/usr/local/share",
				Force: true,
			},
			Expected: `export CDPATH="/usr/local/share:$CDPATH"`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{
			Shell:  tc.Shell,
			CDPath: &context.Path{"/usr/local/share"},
		}
		assert.Equal(t, tc.Expected, tc.CDPath.string(), tc.Case)
	}
}

func TestCDPathRender(t *testing.T) {
	cases := []struct {
		Case           string
		Shell          string
		Expected       string
		CDPaths        CDPaths
		NonEmptyScript bool
	}{
		{
			Case:    "No CDPATH definitions",
			CDPaths: CDPaths{},
			Shell:   BASH,
		},
		{
			Case: "If false",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/share", If: `eq .Shell "fish"`},
			},
			Shell: BASH,
		},
		{
			Case: "If true",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/share", If: `eq .Shell "bash"`},
			},
			Shell:    BASH,
			Expected: `export CDPATH="/usr/share:$CDPATH"`,
		},
		{
			Case: "Single definition with non-empty script",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/share"},
			},
			Shell:          BASH,
			NonEmptyScript: true,
			Expected: `foo

export CDPATH="/usr/share:$CDPATH"`,
		},
		{
			Case: "Two definitions",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/share"},
				&CDPath{Value: "/usr/local/share"},
			},
			Shell: BASH,
			Expected: `export CDPATH="/usr/share:$CDPATH"
export CDPATH="/usr/local/share:$CDPATH"`,
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell, CDPath: &context.Path{}}
		tc.CDPaths.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}

func TestCDPathIfExists(t *testing.T) {
	home := t.TempDir()
	existing := filepath.Join(home, "existing-cd")
	missing := filepath.Join(home, "missing-cd")
	assert.NoError(t, os.MkdirAll(existing, 0o700))

	context.Current = &context.Runtime{
		Shell:  BASH,
		Home:   home,
		CDPath: &context.Path{},
	}

	withExists := &CDPath{
		Value:    Template(existing + "\n" + missing),
		IfExists: true,
	}
	assert.Equal(t, `export CDPATH="`+existing+`:$CDPATH"`, withExists.string(), "ifExists true keeps only existing entries")

	context.Current.CDPath = &context.Path{}
	withoutExists := &CDPath{
		Value:    Template(missing),
		IfExists: false,
	}
	assert.Equal(t, `export CDPATH="`+missing+`:$CDPATH"`, withoutExists.string(), "ifExists false keeps current behavior")
}
