package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type CreatePDBOptions struct {
	genericclioptions.IOStreams

	configFlags *genericclioptions.ConfigFlags
}

func NewCreatePDBOptions(streams genericclioptions.IOStreams) *CreatePDBOptions {
	return &CreatePDBOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

func (o *CreatePDBOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *CreatePDBOptions) Validate() error {
	return nil
}

func (o *CreatePDBOptions) Run() error {
	return nil
}

func NewCmdCreatePDB(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreatePDBOptions(streams)

	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.configFlags.AddFlags(flags)

	return cmd
}
