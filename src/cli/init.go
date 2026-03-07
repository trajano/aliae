package cli

import (
	"os"

	"github.com/jandedobbeleer/aliae/src/appcmd"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	printOutput bool
	ttyOnly     bool

	initCmd = &cobra.Command{
		Use:   "init [bash|zsh|fish|pwsh|powershell|cmd|nu|tcsh|xonsh]",
		Short: "Initialize your shell and config",
		Long: `Initialize your shell and config.

See the documentation to initialize your shell: https://trajano.github.io/aliae/docs/setup/shell.
This is a personal fork. For the official project and docs, see https://aliae.dev.`,
		ValidArgs: []string{
			"bash",
			"zsh",
			"fish",
			"pwsh",
			"powershell",
			"cmd",
			"nu",
			"tcsh",
			"xonsh",
		},
		Args: NoArgsOrOneValidArg,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
			runInit(cmd, args[0])
		},
	}
)

func init() {
	initCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "print the init script")
	initCmd.Flags().BoolVar(&ttyOnly, "tty-only", true, "only print if input is a TTY (default: true)")
	RootCmd.AddCommand(initCmd)
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
