package cli

import (
	"os"

	"github.com/spf13/cobra"
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
This is a personal fork. For the official project and docs, see https://aliae.dev.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if displayVersion {
			cmd.Println(cliVersion)
			return
		}
		_ = cmd.Help()
	},
}

func NewRootCommand(version string) *cobra.Command {
	cliVersion = version
	return RootCmd
}

func Execute(version string) {
	if err := NewRootCommand(version).Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file path")
	RootCmd.Flags().BoolVar(&displayVersion, "version", false, "version")
}
