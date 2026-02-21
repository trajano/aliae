package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	cfg "github.com/jandedobbeleer/aliae/src/config"
)

var (
	printOutput bool
	ttyOnly     bool

	initCmd = &cobra.Command{
		Use:   "init [bash|zsh|fish|pwsh|powershell|cmd|nu|tcsh|xonsh] --tty-only --config ~/.aliae.yaml",
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
			runInit(args[0])
		},
	}
)

func init() {
	initCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "print the init script")
	initCmd.Flags().BoolVar(&ttyOnly, "tty-only", false, "only print if output is a TTY (default: false)")
	_ = initCmd.MarkPersistentFlagRequired("config")
	RootCmd.AddCommand(initCmd)
}

func runInit(shellName string) {
	stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
	stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if shouldSkipInitOutput(ttyOnly, stdinTTY, stdoutTTY) {
		return
	}
	init := cfg.Init(config, shellName, printOutput)
	fmt.Print(init)
}

func shouldSkipInitOutput(ttyOnly, stdinTTY, stdoutTTY bool) bool {
	return ttyOnly && !stdinTTY && !stdoutTTY
}
