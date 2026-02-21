package cli

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
)

var (
	getVerbose bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [shell]",
	Short: "Get a value from aliae",
	Long: `Get a value from aliae.

This command is used to get the value of the following variables:

- shell`,
	ValidArgs: []string{
		"shell",
	},
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		switch args[0] {
		case "shell":
			if getVerbose {
				resolved, trace := shell.NameVerbose()
				for _, line := range trace {
					fmt.Println(line)
				}
				fmt.Printf("shell=%s\n", resolved)
				return
			}
			fmt.Println(shell.Name())
			return
		default:
			_ = cmd.Help()
		}
	},
}

func init() {
	getCmd.Flags().BoolVarP(&getVerbose, "verbose", "v", false, "show shell detection steps")
	RootCmd.AddCommand(getCmd)
}
