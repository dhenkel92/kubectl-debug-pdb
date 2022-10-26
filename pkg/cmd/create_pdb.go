package cmd

import (
	"context"
	"fmt"

	"github.com/dhenkel92/kubectl-utils/pkg/kube"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type CreatePDBOptions struct {
	genericclioptions.IOStreams

	configFlags *genericclioptions.ConfigFlags
	namespace   string
	podName     string

	clients kube.Interface
}

func NewCreatePDBOptions(streams genericclioptions.IOStreams) *CreatePDBOptions {
	return &CreatePDBOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

func (o *CreatePDBOptions) Complete(cmd *cobra.Command, args []string) error {
	o.podName = args[0]

	var err error
	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	cs, err := kube.New(o.configFlags)
	if err != nil {
		return err
	}
	o.clients = cs

	return nil
}

func (o *CreatePDBOptions) Validate() error {
	cs := o.clients.GetClientset()
	_, err := cs.CoreV1().Pods(o.namespace).Get(context.Background(), o.podName, v1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (o *CreatePDBOptions) Run() error {
	return nil
}

func NewCmdCreatePDB(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreatePDBOptions(streams)

	cmd := &cobra.Command{
		Use: "create",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Please pass a pod name")
			}
			if len(args) > 1 {
				return fmt.Errorf("Got too many arguments")
			}
			return nil
		},
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
