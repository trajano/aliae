package cli

import (
	"fmt"

	cfg "github.com/jandedobbeleer/aliae/src/config"
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
		if err := cfg.ValidateConfig(config); err != nil {
			return err
		}

		fmt.Println("configuration is valid")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
