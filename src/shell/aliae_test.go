package shell

import (
	"fmt"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestAliasCommand(t *testing.T) {
	alias := &Alias{Name: "foo", Value: "bar"}
	cases := []struct {
		Case     string
		Shell    string
		Expected string
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Expected: `Set-Alias -Name foo -Value "bar"`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Expected: "macrofile:write(\"foo=bar\", \"\\n\")",
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Expected: `alias foo "bar"`,
		},
		{
			Case:     "NU",
			Shell:    NU,
			Expected: "alias foo = bar",
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: "alias foo 'bar';",
		},
		{
			Case:     "XONSH",
			Shell:    XONSH,
			Expected: "aliases['foo'] = 'bar'",
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Expected: `alias foo="bar"`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Expected: `alias foo="bar"`,
		},
	}

	for _, tc := range cases {
		alias.template = ""
		useRuntime(t, &context.Runtime{Shell: tc.Shell})
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
	}
}

func TestAliasFunction(t *testing.T) {
	cases := []struct {
		Case        string
		Shell       string
		Name        string
		Description string
		Expected    string
	}{
		{
			Case:     "unknown shell",
			Shell:    "unknown",
			Expected: "",
		},
		{
			Case:  "PWSH",
			Shell: PWSH,
			Expected: `function foo() {
    bar
}`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Expected: "",
		},
		{
			Case:  "FISH",
			Shell: FISH,
			Expected: `function foo
    bar
end`,
		},
		{
			Case:        "FISH with description",
			Shell:       FISH,
			Description: "A fish function",
			Expected: `function foo --description "A fish function"
    bar
end`,
		},
		{
			Case:  "NU",
			Shell: NU,
			Expected: `def foo [] {
    bar
}`,
		},
		{
			Case:        "NU with description",
			Shell:       NU,
			Description: "A nu function",
			Expected: `# A nu function
def foo [] {
    bar
}`,
		},
		{
			Case:     "NU",
			Shell:    TCSH,
			Expected: "",
		},
		{
			Case:  "XONSH",
			Shell: XONSH,
			Expected: `@aliases.register("foo")
def __foo():
    bar`,
		},
		{
			Case:  "XONSH - illegal character",
			Name:  "foo-bar",
			Shell: XONSH,
			Expected: `@aliases.register("foo-bar")
def __foobar():
    bar`,
		},
		{
			Case:  "ZSH",
			Shell: ZSH,
			Expected: `foo() {
    bar
}`,
		},
		{
			Case:  "BASH",
			Shell: BASH,
			Expected: `foo() {
    bar
}`,
		},
	}

	for _, tc := range cases {
		alias := &Alias{Name: "foo", Value: "bar", Type: Function}

		if len(tc.Name) > 0 {
			alias.Name = tc.Name
		}
		if len(tc.Description) > 0 {
			alias.Description = tc.Description
		}

		useRuntime(t, &context.Runtime{Shell: tc.Shell})
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
	}
}

func TestAliaeRender(t *testing.T) {
	cases := []struct {
		Case     string
		Expected string
		Aliae    Aliae
	}{
		{
			Case: "Single alias",
			Aliae: Aliae{
				&Alias{Name: "FOO", Value: "bar"},
			},
			Expected: `alias FOO="bar"`,
		},
		{
			Case: "Invalid type",
			Aliae: Aliae{
				&Alias{Name: "FOO", Value: "bar", Type: "invalid"},
			},
		},
		{
			Case: "Double alias",
			Aliae: Aliae{
				&Alias{Name: "FOO", Value: "bar"},
				&Alias{Name: "BAR", Value: "foo"},
			},
			Expected: `alias FOO="bar"
alias BAR="foo"`,
		},
		{
			Case: "Filtered out",
			Aliae: Aliae{
				&Alias{Name: "FOO", Value: "bar", If: `eq .Shell "fish"`},
			},
		},
	}

	for _, tc := range cases {
		ResetRenderOutput()
		useRuntime(t, &context.Runtime{Shell: BASH})
		tc.Aliae.Render()
		assert.Equal(t, tc.Expected, RenderOutputString(), tc.Case)
	}
}

func TestAliasWithTemplate(t *testing.T) {
	cases := []struct {
		Case     string
		Value    Template
		Expected string
	}{
		{
			Case:     "No template",
			Value:    "cd ~",
			Expected: `alias a="cd ~"`,
		},
		{
			Case:     "Home in template",
			Value:    "{{ .Home }}/go/bin/aliae",
			Expected: `alias a="/Users/jan/go/bin/aliae"`,
		},
		{
			Case:     "Advanced template",
			Value:    "{{ .Home }}/go/bin/aliae{{ if eq .OS \"windows\" }}.exe{{ end }}",
			Expected: `alias a="/Users/jan/go/bin/aliae.exe"`,
		},
	}

	for _, tc := range cases {
		alias := &Alias{Name: "a", Value: tc.Value}
		useRuntime(t, &context.Runtime{Shell: BASH, Home: "/Users/jan", OS: context.WINDOWS})
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
	}
}

func TestAliasWithSpacePowerShell(t *testing.T) {
	alias := &Alias{Name: "foo", Value: "bar baz"}
	useRuntime(t, &context.Runtime{Shell: PWSH})
	assert.Equal(t, `function foo() {
	bar baz $args
}`, alias.string())
}

func TestAliasSingleQuoteFish(t *testing.T) {
	alias := &Alias{Name: "foo", Value: "echo 'bar'"}
	useRuntime(t, &context.Runtime{Shell: FISH})
	assert.Equal(t, `alias foo "echo 'bar'"`, alias.string())
}

func TestAliasInlineInterpreters(t *testing.T) {
	interpreterCases := []struct {
		AliasType    Type
		Name         string
		Value        string
		Interpreter  string
		Switch       string
		QuotedForCmd string
	}{
		{
			AliasType:    Python,
			Name:         "pyfoo",
			Value:        "print(123)",
			Interpreter:  "python",
			Switch:       "-c",
			QuotedForCmd: "print(123)",
		},
		{
			AliasType:    Perl,
			Name:         "plfoo",
			Value:        "print qq(123)",
			Interpreter:  "perl",
			Switch:       "-e",
			QuotedForCmd: "print qq(123)",
		},
	}

	shellCases := []struct {
		Case     string
		Shell    string
		Template string
	}{
		{
			Case:  "PWSH",
			Shell: PWSH,
			Template: `function %s() {
    %s %s "%s" $args
}`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Template: `macrofile:write("%s=%s %s \"%s\" $*", "\n")`,
		},
		{
			Case:  "FISH",
			Shell: FISH,
			Template: `function %s
    %s %s "%s" $argv
end`,
		},
		{
			Case:  "NU",
			Shell: NU,
			Template: `def %s [...args] {
    %s %s "%s" ...$args
}`,
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Template: `alias %s '%s %s "%s"';`,
		},
		{
			Case:  "XONSH",
			Shell: XONSH,
			Template: `@aliases.register("%s")
def __%s(args):
    import subprocess
    subprocess.run(["%s", "%s", "%s", *args], check=False)`,
		},
		{
			Case:  "ZSH",
			Shell: ZSH,
			Template: `%s() {
    %s %s "%s" "$@"
}`,
		},
		{
			Case:  "BASH",
			Shell: BASH,
			Template: `%s() {
    %s %s "%s" "$@"
}`,
		},
	}

	for _, interpreter := range interpreterCases {
		for _, shellCase := range shellCases {
			alias := &Alias{Name: interpreter.Name, Value: Template(interpreter.Value), Type: interpreter.AliasType}
			useRuntime(t, &context.Runtime{Shell: shellCase.Shell})

			var expected string
			if shellCase.Shell == XONSH {
				expected = fmt.Sprintf(
					shellCase.Template,
					interpreter.Name,
					interpreter.Name,
					interpreter.Interpreter,
					interpreter.Switch,
					interpreter.Value,
				)
			} else {
				expected = fmt.Sprintf(
					shellCase.Template,
					interpreter.Name,
					interpreter.Interpreter,
					interpreter.Switch,
					interpreter.QuotedForCmd,
				)
			}

			assert.Equal(t, expected, alias.string(), interpreter.Interpreter+" "+shellCase.Case)
		}
	}
}
