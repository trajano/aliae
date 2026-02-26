package shell

import (
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestScriptRender(t *testing.T) {
	cases := []struct {
		Case           string
		Expected       string
		Scripts        Scripts
		NonEmptyScript bool
	}{
		{
			Case:    "No content",
			Scripts: Scripts{},
		},
		{
			Case: "Simple script",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			Expected: "foo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Non-Empty",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			NonEmptyScript: true,
			Expected:       "foo\n\nfoo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Multiple scripts",
			Scripts: Scripts{
				{
					Value: "foo",
				},
				{
					Value: "bar",
				},
			},
			Expected: "foo\nbar",
		},
		{
			Case: "Python script type",
			Scripts: Scripts{
				{
					Type:  PythonScript,
					Value: "print(123)",
				},
			},
			Expected: `python -c "print(123)"`,
		},
		{
			Case: "Perl script type",
			Scripts: Scripts{
				{
					Type:  PerlScript,
					Value: "print qq(123)",
				},
			},
			Expected: `perl -e "print qq(123)"`,
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: PWSH}
		tc.Scripts.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}

func TestScriptRenderCmdPythonAndPerl(t *testing.T) {
	cases := []struct {
		Case     string
		Script   *Script
		Expected string
	}{
		{
			Case: "Python",
			Script: &Script{
				Type:  PythonScript,
				Value: "print(123)",
			},
			Expected: `os.execute("python -c \"print(123)\"")`,
		},
		{
			Case: "Perl",
			Script: &Script{
				Type:  PerlScript,
				Value: "print qq(123)",
			},
			Expected: `os.execute("perl -e \"print qq(123)\"")`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: CMD}
		assert.Equal(t, tc.Expected, tc.Script.String(), tc.Case)
	}
}
