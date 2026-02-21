package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	config         string
	displayVersion bool

	// Version number of aliae
	cliVersion string
)

var RootCmd = &cobra.Command{
	Use:   "aliae",
	Short: "aliae is a tool to do cross platform shell management",
	Long: `aliae is a tool to do cross platform shell management.
It can use the same configuration everywhere to offer a consistent
experience, regardless of where you are. For a detailed guide
on getting started, have a look at the docs at https://trajano.github.io/aliae.
This is a personal fork. For the official project and docs, see https://aliae.dev.
TTY check: stdin=<dynamic> stdout=<dynamic>`,
	Run: func(cmd *cobra.Command, _ []string) {
		if displayVersion {
			fmt.Println(cliVersion)
			return
		}
		stdinTTY := term.IsTerminal(int(os.Stdin.Fd()))
		stdoutTTY := term.IsTerminal(int(os.Stdout.Fd()))
		cmd.Long = fmt.Sprintf(`aliae is a tool to do cross platform shell management.
It can use the same configuration everywhere to offer a consistent
experience, regardless of where you are. For a detailed guide
on getting started, have a look at the docs at https://trajano.github.io/aliae.
This is a personal fork. For the official project and docs, see https://aliae.dev.
TTY check: stdin=%t stdout=%t`, stdinTTY, stdoutTTY)
		_ = cmd.Help()
	},
}

func Execute(version string) {
	cliVersion = version
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file path")
	RootCmd.Flags().BoolVar(&displayVersion, "version", false, "version")
}
