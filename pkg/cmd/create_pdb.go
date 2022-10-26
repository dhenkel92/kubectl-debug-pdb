package cmd

import (
	"fmt"

	"github.com/dhenkel92/kubectl-utils/pkg/kube"
	"github.com/dhenkel92/kubectl-utils/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

type CreatePDBOptions struct {
	genericclioptions.IOStreams

	configFlags *genericclioptions.ConfigFlags
	namespace   string

	workloadName string
	workload     *resource.Info

	clients kube.Interface
}

func NewCreatePDBOptions(streams genericclioptions.IOStreams) *CreatePDBOptions {
	return &CreatePDBOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

func (o *CreatePDBOptions) Complete(cmd *cobra.Command, args []string) error {
	o.workloadName = args[0]

	var err error
	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.clients, err = kube.New(o.configFlags)
	if err != nil {
		return err
	}

	r := o.clients.NewBuilder().
		Unstructured().
		ContinueOnError().
		NamespaceParam(o.namespace).DefaultNamespace().AllNamespaces(false).
		ResourceTypeOrNameArgs(true, o.workloadName).
		Flatten().
		Do()

	infos, err := r.Infos()
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("cannot find workload %s", o.workloadName)
	}
	o.workload = infos[0]

	return nil
}

func (o *CreatePDBOptions) Validate() error {
	if o.workload == nil {
		return fmt.Errorf("cannot find workload '%s'", o.workloadName)
	}
	if !isSupportedWorkload(o.workload) {
		return fmt.Errorf("workload is not supported")
	}

	return nil
}

func isSupportedWorkload(info *resource.Info) bool {
	switch info.Mapping.GroupVersionKind.Kind {
	case "Deployment", "ReplicaSet", "StatefulSet", "Pod":
		return true
	}
	return false
}

func (o *CreatePDBOptions) Run() error {
	workload := o.workload
	unstructured, ok := workload.Object.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("cannot convert workload to unstructured")
	}

	ls, err := utils.LabelsFromWorkload(unstructured)
	if err != nil {
		return err
	}

	pdb := kube.NewPDB(ls)
	fmt.Println(pdb)
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
