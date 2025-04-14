package cmd

import (
	"fmt"
	"io"
	"os"

	fiddle "github.com/richardmarshall/vclfiddle"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	Writer   io.Writer
	FiddleId string
	Client   *fiddle.Client
}

func NewGetCommand(w io.Writer) *cobra.Command {
	o := NewGetOptions(w)
	c := &cobra.Command{
		Use:   "get",
		Short: "lookup a fiddle by id",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(args); err != nil {
				return err
			}
			return o.Run()
		},
	}

	return c
}

func NewGetOptions(w io.Writer) *GetOptions {
	if w == nil {
		w = os.Stdout
	}
	return &GetOptions{Writer: w}
}

func (o *GetOptions) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("provide single fiddle id")
	}

	o.FiddleId = args[0]
	o.Client = fiddle.NewClient()
	return nil
}

func (o *GetOptions) Run() error {
	f, valid, lints, err := o.Client.Get(o.FiddleId)
	if err != nil {
		return err
	}
	if !valid {
		fmt.Fprintf(o.Writer, "%#v\n", lints)
	}
	fmt.Fprintln(o.Writer, fiddle.PrettyPrint(f))
	return nil
}
