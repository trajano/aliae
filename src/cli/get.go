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
	normalizeEffectiveScriptWeights(aliae)

	data, err := yaml.Marshal(aliae)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func normalizeEffectiveScriptWeights(aliae *cfg.Aliae) {
	if aliae == nil {
		return
	}

	for _, script := range aliae.Scripts {
		if script == nil {
			continue
		}

		if script.Weight <= 0 {
			script.Weight = 1
		}
	}
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
