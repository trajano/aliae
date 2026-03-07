package cli

import (
	"github.com/jandedobbeleer/aliae/src/appcmd"
	"github.com/spf13/cobra"
)

var stateCmd = newStateCommand()
var stateListCmd = newStateListCommand()
var stateClearCmd = newStateClearCommand()

func newStateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "state [list|clear]",
		Short: "Inspect and manage aliae script state files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStateList(cmd)
		},
	}
}

func newStateListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List referenced aliae script state files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStateList(cmd)
		},
	}
}

func newStateClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Remove referenced aliae script state files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStateClear(cmd)
		},
	}
}

func registerStateCommand(root *cobra.Command) {
	stateCmd = newStateCommand()
	stateListCmd = newStateListCommand()
	stateClearCmd = newStateClearCommand()
	stateCmd.AddCommand(stateListCmd)
	stateCmd.AddCommand(stateClearCmd)
	root.AddCommand(stateCmd)
}

func runStateList(cmd *cobra.Command) error {
	return appcmd.StateListCommand{
		ConfigPath: config,
		Out:        cmd.OutOrStdout(),
	}.Execute()
}

func runStateClear(cmd *cobra.Command) error {
	return appcmd.StateClearCommand{
		ConfigPath: config,
		Out:        cmd.OutOrStdout(),
	}.Execute()
}
