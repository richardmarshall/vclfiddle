package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewDefaultCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "fiddle",
		Short: "",
		Long:  ``,
	}

	c.AddCommand(
		NewGetCommand(os.Stdout),
		NewExecuteCommand(os.Stdout),
	)

	return c
}
