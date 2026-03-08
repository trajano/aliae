package cli

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/appcmd"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = newGetCommand()

func newGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [shell|config|variables|benchmark [shell]]",
		Short: "Get a value from aliae",
		Long: `Get a value from aliae.

This command is used to get the value of the following variables:

- shell
- config
- variables
- benchmark [shell]`,
		ValidArgs: []string{
			"shell",
			"config",
			"variables",
			"benchmark",
		},
		Args: validateGetArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "shell":
				return appcmd.GetShellCommand{Out: cmd.OutOrStdout()}.Execute()
			case "config":
				return appcmd.GetResolvedConfigCommand{
					ConfigPath: config,
					Out:        cmd.OutOrStdout(),
				}.Execute()
			case "variables":
				return appcmd.GetVariablesCommand{
					ConfigPath: config,
					Out:        cmd.OutOrStdout(),
				}.Execute()
			case "benchmark":
				benchmarkShell := ""
				if len(args) == 2 {
					benchmarkShell = args[1]
				}
				return appcmd.BenchmarkCommand{
					ConfigPath:     config,
					BenchmarkShell: benchmarkShell,
					NoCache:        benchmarkNoCache,
					Out:            cmd.OutOrStdout(),
				}.Execute()
			}

			return nil
		},
	}
}

var benchmarkNoCache bool

func registerGetCommand(root *cobra.Command) {
	benchmarkNoCache = false
	getCmd = newGetCommand()
	getCmd.Flags().BoolVar(&benchmarkNoCache, "no-cache", false, "disable config cache while running benchmark")
	root.AddCommand(getCmd)
}

func validateGetArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 2 {
		return fmt.Errorf("accepts 1 or 2 args, received %d", len(args))
	}

	switch args[0] {
	case "shell", "config", "variables":
		if len(args) != 1 {
			return fmt.Errorf("%s does not accept extra arguments", args[0])
		}
		return nil
	case "benchmark":
		if len(args) == 1 {
			return nil
		}

		if !isSupportedBenchmarkShell(args[1]) {
			return fmt.Errorf("unsupported shell for benchmark: %s", args[1])
		}
		return nil
	default:
		return fmt.Errorf("invalid argument %q, expected shell, config, variables, or benchmark", args[0])
	}
}

func isSupportedBenchmarkShell(shellName string) bool {
	switch shellName {
	case shell.BASH, shell.ZSH, shell.PWSH, shell.POWERSHELL, shell.FISH, shell.NU, shell.TCSH, shell.XONSH, shell.CMD:
		return true
	default:
		return false
	}
}

func renderResolvedConfigYAML(configPath string) (string, error) {
	return appcmd.RenderResolvedConfigYAML(configPath)
}
