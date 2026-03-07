package shell

type Aliae []*Alias

type Alias struct {
	Name        string   `yaml:"name"`
	Value       Template `yaml:"value"`
	Type        Type     `yaml:"type"`
	If          If       `yaml:"if"`
	Description string   `yaml:"description"`
	Option      Option   `yaml:"option"`
	Scope       Option   `yaml:"scope"`
	template    string
	Force       bool `yaml:"force"`
}

type Option string

type Type string

const (
	Command  Type = "command"
	Function Type = "function"
	Git      Type = "git"
	Python   Type = "python"
	Perl     Type = "perl"
)

func (a *Alias) string() string {
	if len(a.Type) == 0 {
		a.Type = Command
	}

	if a.Type == Git {
		return a.git()
	}

	a.Value = a.Value.Parse()

	return formatStrategy().FormatAlias(a)
}

func (a *Alias) render() string {
	script, err := parse(a.template, a)
	if err != nil {
		return err.Error()
	}

	return script
}

func (a Aliae) Render() {
	if len(a) == 0 {
		return
	}

	first := true
	strategy := formatStrategy()
	wrotePrelude := false
	for _, alias := range a {
		if alias.If.Ignore() {
			continue
		}

		script := alias.string()
		if len(script) == 0 {
			advanceAutoProgress(1)
			continue
		}

		if first && dotFileHasRenderableContent() {
			writeRenderOutput("\n\n")
		}

		if first {
			prelude := strategy.FormatAliasScriptPrelude()
			if prelude != "" {
				writeRenderOutput(prelude)
				wrotePrelude = true
			}
		}

		if !first {
			if !dotFileEndsWithNewline() {
				writeRenderOutput("\n")
			}
		}

		writeRenderOutput(script)

		first = false
		advanceAutoProgress(1)
	}

	if wrotePrelude {
		writeRenderOutput(strategy.FormatAliasScriptPostlude())
	}
}
