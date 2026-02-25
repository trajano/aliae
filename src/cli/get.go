package cli

import (
	"fmt"
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
			fmt.Println(shell.Name())
			return nil
		case "config":
			output, err := renderResolvedConfigYAML(config)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		case "variables":
			return printVariableDiagnostics()
		default:
			_ = cmd.Help()
			return nil
		}
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

	data, err := yaml.Marshal(aliae)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// printVariableDiagnostics prints runtime and template diagnostics for 'aliae get variables'.
func printVariableDiagnostics() error {
	shellName, trace := shell.NameVerbose()
	context.Init(shellName)

	configPath, configDir := cfg.ResolveTemplateContext(config)
	context.Current.ConfigPath = configPath
	context.Current.ConfigDir = configDir

	stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
	stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))

	fmt.Fprintln(os.Stdout, "aliae get variables")
	fmt.Fprintf(os.Stdout, "tty.stdin=%t\n", stdinTTY)
	fmt.Fprintf(os.Stdout, "tty.stdout=%t\n", stdoutTTY)
	fmt.Fprintf(os.Stdout, "template.Shell=%s\n", context.Current.Shell)
	fmt.Fprintf(os.Stdout, "template.OS=%s\n", context.Current.OS)
	fmt.Fprintf(os.Stdout, "template.WSL=%t\n", context.Current.WSL)
	fmt.Fprintf(os.Stdout, "template.Hostname=%s\n", context.Current.Hostname)
	fmt.Fprintf(os.Stdout, "template.Home=%s\n", context.Current.Home)
	fmt.Fprintf(os.Stdout, "template.Arch=%s\n", context.Current.Arch)
	fmt.Fprintf(os.Stdout, "template.ConfigPath=%s\n", context.Current.ConfigPath)
	fmt.Fprintf(os.Stdout, "template.ConfigDir=%s\n", context.Current.ConfigDir)
	for _, line := range trace {
		fmt.Fprintf(os.Stdout, "shell.trace=%s\n", line)
	}

	return nil
}
