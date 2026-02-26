package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [shell|config|variables]",
	Short: "Get a value from aliae",
	Long: `Get a value from aliae.

This command is used to get the value of the following variables:

- shell
- config
- variables`,
	ValidArgs: []string{
		"shell",
		"config",
		"variables",
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "shell":
			cmd.Println(shell.Name())
			return nil
		case "config":
			output, err := renderResolvedConfigYAML(config)
			if err != nil {
				return err
			}
			cmd.Println(output)
			return nil
		case "variables":
			return printVariableDiagnostics(cmd.OutOrStdout())
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func renderResolvedConfigYAML(configPath string) (string, error) {
	aliae, err := cfg.LoadConfig(configPath)
	if err != nil {
		return "", err
	}

	output := resolvedConfigOutput{
		Alias:    toAliasOutput(aliae.Aliae),
		Env:      toEnvOutput(aliae.Envs),
		Path:     toPathOutput(aliae.Paths),
		CDPath:   toCDPathOutput(aliae.CDPaths),
		Script:   toScriptOutput(aliae.Scripts),
		Link:     toLinkOutput(aliae.Links),
		Progress: aliae.Progress,
	}
	if aliae.StatTimeout > 0 {
		output.StatTimeout = aliae.StatTimeout.String()
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type resolvedConfigOutput struct {
	Alias       []aliasOutput  `yaml:"alias,omitempty"`
	Env         []envOutput    `yaml:"env,omitempty"`
	Path        []pathOutput   `yaml:"path,omitempty"`
	CDPath      []pathOutput   `yaml:"cdpath,omitempty"`
	Script      []scriptOutput `yaml:"script,omitempty"`
	Link        []linkOutput   `yaml:"link,omitempty"`
	Progress    cfg.Progress   `yaml:"progress,omitempty"`
	StatTimeout string         `yaml:"stat_timeout,omitempty"`
}

type aliasOutput struct {
	Name        string         `yaml:"name"`
	Value       shell.Template `yaml:"value"`
	Type        shell.Type     `yaml:"type,omitempty"`
	If          shell.If       `yaml:"if,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Option      shell.Option   `yaml:"option,omitempty"`
	Scope       shell.Option   `yaml:"scope,omitempty"`
	Force       bool           `yaml:"force,omitempty"`
}

type envOutput struct {
	Name      string         `yaml:"name"`
	Value     any            `yaml:"value"`
	Delimiter shell.Template `yaml:"delimiter,omitempty"`
	If        shell.If       `yaml:"if,omitempty"`
	Type      shell.EnvType  `yaml:"type,omitempty"`
	IsPath    bool           `yaml:"isPath,omitempty"`
	IfExists  bool           `yaml:"ifExists,omitempty"`
	Persist   bool           `yaml:"persist,omitempty"`
}

type pathOutput struct {
	Value    shell.Template `yaml:"value"`
	If       shell.If       `yaml:"if,omitempty"`
	Persist  bool           `yaml:"persist,omitempty"`
	Force    bool           `yaml:"force,omitempty"`
	IfExists bool           `yaml:"ifExists,omitempty"`
}

type scriptOutput struct {
	Value  shell.Template `yaml:"value"`
	If     shell.If       `yaml:"if,omitempty"`
	Weight float64        `yaml:"weight,omitempty"`
}

type linkOutput struct {
	Name   shell.Template `yaml:"name"`
	Target shell.Template `yaml:"target"`
	If     shell.If       `yaml:"if,omitempty"`
	MkDir  bool           `yaml:"mkdir,omitempty"`
}

func toAliasOutput(items shell.Aliae) []aliasOutput {
	output := make([]aliasOutput, 0, len(items))
	for _, alias := range items {
		if alias == nil {
			continue
		}

		output = append(output, aliasOutput{
			Name:        alias.Name,
			Value:       alias.Value,
			Type:        alias.Type,
			If:          alias.If,
			Description: alias.Description,
			Option:      alias.Option,
			Scope:       alias.Scope,
			Force:       alias.Force,
		})
	}

	return output
}

func toEnvOutput(items shell.Envs) []envOutput {
	output := make([]envOutput, 0, len(items))
	for _, env := range items {
		if env == nil {
			continue
		}

		output = append(output, envOutput{
			Name:      env.Name,
			Value:     env.Value,
			Delimiter: env.Delimiter,
			If:        env.If,
			Type:      env.Type,
			IsPath:    env.IsPath,
			IfExists:  env.IfExists,
			Persist:   env.Persist,
		})
	}

	return output
}

func toPathOutput(items shell.Paths) []pathOutput {
	output := make([]pathOutput, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		output = append(output, pathOutput{
			Value:    item.Value,
			If:       item.If,
			Persist:  item.Persist,
			Force:    item.Force,
			IfExists: item.IfExists,
		})
	}

	return output
}

func toCDPathOutput(items shell.CDPaths) []pathOutput {
	output := make([]pathOutput, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		output = append(output, pathOutput{
			Value:    item.Value,
			If:       item.If,
			Force:    item.Force,
			IfExists: item.IfExists,
		})
	}

	return output
}

func toScriptOutput(items shell.Scripts) []scriptOutput {
	output := make([]scriptOutput, 0, len(items))
	for _, script := range items {
		if script == nil {
			continue
		}

		entry := scriptOutput{
			Value: script.Value,
			If:    script.If,
		}
		if script.Weight > 0 {
			entry.Weight = script.Weight
		}

		output = append(output, entry)
	}

	return output
}

func toLinkOutput(items shell.Links) []linkOutput {
	output := make([]linkOutput, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		output = append(output, linkOutput{
			Name:   item.Name,
			Target: item.Target,
			If:     item.If,
			MkDir:  item.MkDir,
		})
	}

	return output
}

// printVariableDiagnostics prints runtime and template diagnostics for 'aliae get variables'.
func printVariableDiagnostics(out io.Writer) error {
	shellName, trace := shell.NameVerbose()
	context.Init(shellName)

	configPath, configDir := cfg.ResolveTemplateContext(config)
	context.Current.ConfigPath = configPath
	context.Current.ConfigDir = configDir

	stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
	stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))

	fmt.Fprintln(out, "aliae get variables")
	fmt.Fprintf(out, "tty.stdin=%t\n", stdinTTY)
	fmt.Fprintf(out, "tty.stdout=%t\n", stdoutTTY)
	fmt.Fprintf(out, "template.Shell=%s\n", context.Current.Shell)
	fmt.Fprintf(out, "template.OS=%s\n", context.Current.OS)
	fmt.Fprintf(out, "template.WSL=%t\n", context.Current.WSL)
	fmt.Fprintf(out, "template.Hostname=%s\n", context.Current.Hostname)
	fmt.Fprintf(out, "template.Home=%s\n", context.Current.Home)
	fmt.Fprintf(out, "template.Arch=%s\n", context.Current.Arch)
	fmt.Fprintf(out, "template.ConfigPath=%s\n", context.Current.ConfigPath)
	fmt.Fprintf(out, "template.ConfigDir=%s\n", context.Current.ConfigDir)
	for _, line := range trace {
		fmt.Fprintf(out, "shell.trace=%s\n", line)
	}

	return nil
}
