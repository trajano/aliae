package cli

import (
	"fmt"

	"github.com/goccy/go-yaml"
	cfg "github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [shell|config]",
	Short: "Get a value from aliae",
	Long: `Get a value from aliae.

This command is used to get the value of the following variables:

- shell
- config`,
	ValidArgs: []string{
		"shell",
		"config",
	},
	Args: cobra.OnlyValidArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "shell":
			fmt.Println(shell.Name())
			return nil
		case "config":
			output, err := renderResolvedConfigYAML(config)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		default:
			_ = cmd.Help()
			return nil
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func renderResolvedConfigYAML(configPath string) (string, error) {
	aliae, err := cfg.LoadConfig(configPath)
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(aliae)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
