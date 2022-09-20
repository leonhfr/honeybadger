package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// optionsCmd represents the options command.
// It returns the list of available options.
var optionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Lists the available options",
	Long: `Options lists the available options in UCI format.

Fields:
  name        name
  type        type either spin (integer) or combo (enum)
  default     default value
  min         minimum value
  max         maximum value
  var         one of enum possible value`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		e := newEngine(ctx)

		options := e.Options()
		for _, option := range options {
			fmt.Println(option)
		}

		return nil
	},
}
