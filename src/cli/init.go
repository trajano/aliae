package cli

import (
	"os"

	"github.com/jandedobbeleer/aliae/src/appcmd"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	printOutput bool
	ttyOnly     bool

	initCmd = newInitCommand()
)

var initValidArgs = []string{
	shell.BASH,
	shell.ZSH,
	shell.FISH,
	shell.PWSH,
	shell.POWERSHELL,
	shell.CMD,
	shell.NU,
	shell.TCSH,
	shell.XONSH,
}

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init [bash|zsh|fish|pwsh|powershell|cmd|nu|tcsh|xonsh]",
		Short: "Initialize your shell and config",
		Long: `Initialize your shell and config.

See the documentation to initialize your shell: https://trajano.github.io/aliae/docs/setup/shell.
This is a personal fork. For the official project and docs, see https://aliae.dev.`,
		ValidArgs: initValidArgs,
		Args:      NoArgsOrOneValidArg,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
			runInit(cmd, args[0])
		},
	}
}

func registerInitCommand(root *cobra.Command) {
	printOutput = false
	ttyOnly = false
	initCmd = newInitCommand()
	initCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "print the init script")
	initCmd.Flags().BoolVar(&ttyOnly, "tty-only", true, "only print if input is a TTY (default: true)")
	root.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, shellName string) {
	_ = appcmd.InitCommand{
		ConfigPath: config,
		Shell:      shellName,
		Print:      printOutput,
		TTYOnly:    ttyOnly,
		StdinTTY:   term.IsTerminal(int(os.Stdin.Fd())),
		Out:        cmd.OutOrStdout(),
		Err:        cmd.ErrOrStderr(),
	}.Execute()
}

func shouldSkipInitOutput(ttyOnly, stdinTTY bool) bool {
	return ttyOnly && !stdinTTY
}
