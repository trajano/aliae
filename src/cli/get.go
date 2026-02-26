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
			fmt.Fprintln(cmd.OutOrStdout(), shell.Name())
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

	output := yaml.MapSlice{}
	if aliases := toAliasOutput(aliae.Aliae); len(aliases) > 0 {
		output = append(output, yaml.MapItem{Key: "alias", Value: aliases})
	}
	if envs := toEnvOutput(aliae.Envs); len(envs) > 0 {
		output = append(output, yaml.MapItem{Key: "env", Value: envs})
	}
	if paths := toPathOutput(aliae.Paths); len(paths) > 0 {
		output = append(output, yaml.MapItem{Key: "path", Value: paths})
	}
	if cdpaths := toCDPathOutput(aliae.CDPaths); len(cdpaths) > 0 {
		output = append(output, yaml.MapItem{Key: "cdpath", Value: cdpaths})
	}
	if scripts := toScriptOutput(aliae.Scripts); len(scripts) > 0 {
		output = append(output, yaml.MapItem{Key: "script", Value: scripts})
	}
	if links := toLinkOutput(aliae.Links); len(links) > 0 {
		output = append(output, yaml.MapItem{Key: "link", Value: links})
	}
	if aliae.Progress.Enabled {
		output = append(output, yaml.MapItem{Key: "progress", Value: aliae.Progress})
	}
	if aliae.StatTimeout > 0 {
		output = append(output, yaml.MapItem{Key: "stat_timeout", Value: aliae.StatTimeout.String()})
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func toAliasOutput(items shell.Aliae) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, alias := range items {
		if alias == nil {
			continue
		}

		item := yaml.MapSlice{{Key: "name", Value: alias.Name}}
		if len(alias.Type) > 0 {
			item = append(item, yaml.MapItem{Key: "type", Value: alias.Type})
		}
		item = append(item, yaml.MapItem{Key: "value", Value: alias.Value})
		if len(alias.Description) > 0 {
			item = append(item, yaml.MapItem{Key: "description", Value: alias.Description})
		}
		if alias.Force {
			item = append(item, yaml.MapItem{Key: "force", Value: true})
		}
		if len(alias.If) > 0 {
			item = append(item, yaml.MapItem{Key: "if", Value: alias.If})
		}
		if len(alias.Option) > 0 {
			item = append(item, yaml.MapItem{Key: "option", Value: alias.Option})
		}
		if len(alias.Scope) > 0 {
			item = append(item, yaml.MapItem{Key: "scope", Value: alias.Scope})
		}

		output = append(output, item)
	}

	return output
}

func toEnvOutput(items shell.Envs) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, env := range items {
		if env == nil {
			continue
		}

		item := yaml.MapSlice{{Key: "name", Value: env.Name}}
		if len(env.Type) > 0 {
			item = append(item, yaml.MapItem{Key: "type", Value: env.Type})
		}
		item = append(item, yaml.MapItem{Key: "value", Value: env.Value})
		if len(env.Delimiter) > 0 {
			item = append(item, yaml.MapItem{Key: "delimiter", Value: env.Delimiter})
		}
		if len(env.If) > 0 {
			item = append(item, yaml.MapItem{Key: "if", Value: env.If})
		}
		if env.IfExists {
			item = append(item, yaml.MapItem{Key: "ifExists", Value: true})
		}
		if env.IsPath {
			item = append(item, yaml.MapItem{Key: "isPath", Value: true})
		}
		if env.Persist {
			item = append(item, yaml.MapItem{Key: "persist", Value: true})
		}

		output = append(output, item)
	}

	return output
}

func toPathOutput(items shell.Paths) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: "value", Value: item.Value}}
		if item.Force {
			entry = append(entry, yaml.MapItem{Key: "force", Value: true})
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.IfExists {
			entry = append(entry, yaml.MapItem{Key: "ifExists", Value: true})
		}
		if item.Persist {
			entry = append(entry, yaml.MapItem{Key: "persist", Value: true})
		}

		output = append(output, entry)
	}

	return output
}

func toCDPathOutput(items shell.CDPaths) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: "value", Value: item.Value}}
		if item.Force {
			entry = append(entry, yaml.MapItem{Key: "force", Value: true})
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.IfExists {
			entry = append(entry, yaml.MapItem{Key: "ifExists", Value: true})
		}

		output = append(output, entry)
	}

	return output
}

func toScriptOutput(items shell.Scripts) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, script := range items {
		if script == nil {
			continue
		}

		entry := yaml.MapSlice{{Key: "value", Value: script.Value}}
		if len(script.Type) > 0 {
			entry = append(entry, yaml.MapItem{Key: "type", Value: script.Type})
		}
		if len(script.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: script.If})
		}
		if script.Weight > 0 {
			entry = append(entry, yaml.MapItem{Key: "weight", Value: script.Weight})
		}
		if len(script.State.File) > 0 {
			stateEntry := yaml.MapSlice{{Key: "file", Value: script.State.File}}
			if len(script.State.RunEvery) > 0 {
				stateEntry = append(stateEntry, yaml.MapItem{Key: "runEvery", Value: script.State.RunEvery})
			}
			if len(script.State.Format) > 0 {
				stateEntry = append(stateEntry, yaml.MapItem{Key: "format", Value: script.State.Format})
			}
			entry = append(entry, yaml.MapItem{Key: "state", Value: stateEntry})
		}

		output = append(output, entry)
	}

	return output
}

func toLinkOutput(items shell.Links) []yaml.MapSlice {
	output := make([]yaml.MapSlice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := yaml.MapSlice{
			{Key: "name", Value: item.Name},
			{Key: "target", Value: item.Target},
		}
		if len(item.If) > 0 {
			entry = append(entry, yaml.MapItem{Key: "if", Value: item.If})
		}
		if item.MkDir {
			entry = append(entry, yaml.MapItem{Key: "mkdir", Value: true})
		}

		output = append(output, entry)
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
