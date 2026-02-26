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

	output := map[string]any{}
	if aliases := toAliasOutput(aliae.Aliae); len(aliases) > 0 {
		output["alias"] = aliases
	}
	if envs := toEnvOutput(aliae.Envs); len(envs) > 0 {
		output["env"] = envs
	}
	if paths := toPathOutput(aliae.Paths); len(paths) > 0 {
		output["path"] = paths
	}
	if cdpaths := toCDPathOutput(aliae.CDPaths); len(cdpaths) > 0 {
		output["cdpath"] = cdpaths
	}
	if scripts := toScriptOutput(aliae.Scripts); len(scripts) > 0 {
		output["script"] = scripts
	}
	if links := toLinkOutput(aliae.Links); len(links) > 0 {
		output["link"] = links
	}
	if aliae.Progress.Enabled {
		output["progress"] = aliae.Progress
	}
	if aliae.StatTimeout > 0 {
		output["stat_timeout"] = aliae.StatTimeout.String()
	}

	data, err := yaml.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func toAliasOutput(items shell.Aliae) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, alias := range items {
		if alias == nil {
			continue
		}

		item := map[string]any{
			"name":  alias.Name,
			"value": alias.Value,
		}
		if len(alias.Type) > 0 {
			item["type"] = alias.Type
		}
		if len(alias.If) > 0 {
			item["if"] = alias.If
		}
		if len(alias.Description) > 0 {
			item["description"] = alias.Description
		}
		if len(alias.Option) > 0 {
			item["option"] = alias.Option
		}
		if len(alias.Scope) > 0 {
			item["scope"] = alias.Scope
		}
		if alias.Force {
			item["force"] = true
		}

		output = append(output, item)
	}

	return output
}

func toEnvOutput(items shell.Envs) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, env := range items {
		if env == nil {
			continue
		}

		item := map[string]any{
			"name":  env.Name,
			"value": env.Value,
		}
		if len(env.Delimiter) > 0 {
			item["delimiter"] = env.Delimiter
		}
		if len(env.If) > 0 {
			item["if"] = env.If
		}
		if len(env.Type) > 0 {
			item["type"] = env.Type
		}
		if env.IsPath {
			item["isPath"] = true
		}
		if env.IfExists {
			item["ifExists"] = true
		}
		if env.Persist {
			item["persist"] = true
		}

		output = append(output, item)
	}

	return output
}

func toPathOutput(items shell.Paths) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := map[string]any{
			"value": item.Value,
		}
		if len(item.If) > 0 {
			entry["if"] = item.If
		}
		if item.Persist {
			entry["persist"] = true
		}
		if item.Force {
			entry["force"] = true
		}
		if item.IfExists {
			entry["ifExists"] = true
		}

		output = append(output, entry)
	}

	return output
}

func toCDPathOutput(items shell.CDPaths) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := map[string]any{
			"value": item.Value,
		}
		if len(item.If) > 0 {
			entry["if"] = item.If
		}
		if item.Force {
			entry["force"] = true
		}
		if item.IfExists {
			entry["ifExists"] = true
		}

		output = append(output, entry)
	}

	return output
}

func toScriptOutput(items shell.Scripts) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, script := range items {
		if script == nil {
			continue
		}

		entry := map[string]any{
			"value": script.Value,
		}
		if len(script.If) > 0 {
			entry["if"] = script.If
		}
		if script.Weight > 0 {
			entry["weight"] = script.Weight
		}

		output = append(output, entry)
	}

	return output
}

func toLinkOutput(items shell.Links) []map[string]any {
	output := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		entry := map[string]any{
			"name":   item.Name,
			"target": item.Target,
		}
		if len(item.If) > 0 {
			entry["if"] = item.If
		}
		if item.MkDir {
			entry["mkdir"] = true
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
