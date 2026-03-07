package cli

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/appcmd"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:          "validate",
	Short:        "Validate config against the schema",
	SilenceUsage: true,
	Long: `Validate config against the schema.

This command validates the resolved configuration after extends and include directives are applied.`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := (appcmd.ValidateCommand{ConfigPath: config}).Execute(); err != nil {
			return err
		}

		fmt.Println("configuration is valid")
		return nil
	},
}

func registerValidateCommand(root *cobra.Command) {
	root.AddCommand(validateCmd)
}
