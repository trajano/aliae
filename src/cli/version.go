package cli

import "github.com/spf13/cobra"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  "Print the version number of aliae.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println(cliVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
