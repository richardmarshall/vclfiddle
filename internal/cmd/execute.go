package cmd

import (
	"fmt"
	"io"
	"os"

	fiddle "github.com/richardmarshall/vclfiddle"
	"github.com/spf13/cobra"
)

type ExecuteOptions struct {
	Writer   io.Writer
	FiddleId string
	CacheID  int
	Client   *fiddle.Client
}

func NewExecuteCommand(w io.Writer) *cobra.Command {
	o := NewExecuteOptions(w)
	c := &cobra.Command{
		Use:   "execute",
		Short: "execute a fiddle",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(args); err != nil {
				return err
			}
			return o.Run()
		},
	}

	c.Flags().IntVar(&o.CacheID, "cache-id", 0, "")

	return c
}

func NewExecuteOptions(w io.Writer) *ExecuteOptions {
	if w == nil {
		w = os.Stdout
	}
	return &ExecuteOptions{Writer: w}
}

func (o *ExecuteOptions) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("provide single fiddle id")
	}

	if o.CacheID > 100000 {
		return fmt.Errorf("cache-id must be smaller than 100000")
	}

	o.FiddleId = args[0]
	o.Client = fiddle.NewClient()
	return nil
}

func (o *ExecuteOptions) Run() error {
	f, err := o.Client.Get(o.FiddleId)
	if err != nil {
		return err
	}
	r, err := o.Client.Execute(f, fiddle.ExecuteOptions{CacheID: o.CacheID})
	if err != nil {
		return err
	}
	fmt.Fprintln(o.Writer, fiddle.PrettyPrint(r))
	return nil
}
