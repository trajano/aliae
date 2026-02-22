package cli

import (
	"fmt"
	"os"

	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Print runtime and template diagnostics",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		shellName, trace := shell.NameVerbose()
		context.Init(shellName)

		configPath, configDir := cfg.ResolveTemplateContext(config)
		context.Current.ConfigPath = configPath
		context.Current.ConfigDir = configDir
		context.Current.AliaeConfig = configPath
		context.Current.AliaeConfigDir = configDir

		stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
		stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))

		fmt.Fprintln(os.Stdout, "aliae debug")
		fmt.Fprintf(os.Stdout, "tty.stdin=%t\n", stdinTTY)
		fmt.Fprintf(os.Stdout, "tty.stdout=%t\n", stdoutTTY)
		fmt.Fprintf(os.Stdout, "template.Shell=%s\n", context.Current.Shell)
		fmt.Fprintf(os.Stdout, "template.OS=%s\n", context.Current.OS)
		fmt.Fprintf(os.Stdout, "template.Hostname=%s\n", context.Current.Hostname)
		fmt.Fprintf(os.Stdout, "template.Home=%s\n", context.Current.Home)
		fmt.Fprintf(os.Stdout, "template.Arch=%s\n", context.Current.Arch)
		fmt.Fprintf(os.Stdout, "template.ConfigPath=%s\n", context.Current.ConfigPath)
		fmt.Fprintf(os.Stdout, "template.ConfigDir=%s\n", context.Current.ConfigDir)
		fmt.Fprintf(os.Stdout, "template.AliaeConfig=%s\n", context.Current.AliaeConfig)
		fmt.Fprintf(os.Stdout, "template.AliaeConfigDir=%s\n", context.Current.AliaeConfigDir)
		for _, line := range trace {
			fmt.Fprintf(os.Stdout, "shell.trace=%s\n", line)
		}
	},
}

func init() {
	RootCmd.AddCommand(debugCmd)
}
