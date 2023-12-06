package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/dhenkel92/kubectl-pdb/pkg/kube"
	"github.com/dhenkel92/kubectl-pdb/pkg/utils"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	policyv1 "k8s.io/api/policy/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
)

type CreatePDBOptions struct {
	genericclioptions.IOStreams

	configFlags genericclioptions.RESTClientGetter
	namespace   string
	output      string
	dryRun      bool

	minAvail      *intstr.IntOrString
	minAvailStr   string
	maxUnavail    *intstr.IntOrString
	maxUnavailStr string

	workloadName string
	workload     *resource.Info

	clients kube.Interface
}

func NewCreatePDBOptions(streams genericclioptions.IOStreams, conf genericclioptions.RESTClientGetter) *CreatePDBOptions {
	return &CreatePDBOptions{
		IOStreams:   streams,
		configFlags: conf,
	}
}

func (o *CreatePDBOptions) GetPrinter() (printers.ResourcePrinter, error) {
	switch o.output {
	case "json":
		return &printers.JSONPrinter{}, nil
	case "yaml":
		return &printers.YAMLPrinter{}, nil
	}
	return nil, fmt.Errorf("'%s' output is not supported", o.output)
}

func (o *CreatePDBOptions) getWriter() *tabwriter.Writer {
	return printers.GetNewTabWriter(o.Out)
}

func (o *CreatePDBOptions) Complete(cmd *cobra.Command, args []string) error {
	o.output = strings.Trim(o.output, " ")

	var err error
	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	o.clients, err = kube.New(o.configFlags)
	if err != nil {
		return err
	}

	o.workloadName = args[0]
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

	if o.minAvailStr != "" {
		o.minAvail = utils.StrToIntOrString(o.minAvailStr)
	}
	if o.maxUnavailStr != "" {
		o.maxUnavail = utils.StrToIntOrString(o.maxUnavailStr)
	}
	if o.maxUnavail == nil && o.minAvail == nil {
		o.maxUnavail = utils.StrToIntOrString("1")
	}

	return nil
}

func (o *CreatePDBOptions) Validate() error {
	if o.workload == nil {
		return fmt.Errorf("cannot find workload '%s'", o.workloadName)
	}
	if !isSupportedWorkload(o.workload) {
		return fmt.Errorf("workload is not supported")
	}
	if o.minAvail != nil && o.maxUnavail != nil {
		return fmt.Errorf("you can either set min-available or max-available")
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

	pdb := kube.NewPDB(unstructured, ls, o.minAvail, o.maxUnavail)
	if !o.dryRun {
		pdb, err = o.clients.GetClientset().PolicyV1().PodDisruptionBudgets(o.namespace).Create(context.Background(), pdb, v1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	// we have to set the group version kind, so that the printer can do it's job
	// missing apiVersion or kind; try GetObjectKind().SetGroupVersionKind() if you know the type
	pdb.GetObjectKind().SetGroupVersionKind(policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget"))
	if o.output == "name" {
		fmt.Fprintf(o.Out, "poddisruptionbudget/%s created\n", pdb.GetName())
	} else {
		printer, err := o.GetPrinter()
		if err != nil {
			return err
		}
		return printer.PrintObj(pdb, o.getWriter())
	}

	return nil
}

func NewCmdCreatePDB(streams genericclioptions.IOStreams, conf genericclioptions.RESTClientGetter) *cobra.Command {
	o := NewCreatePDBOptions(streams, conf)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new PDB for a given workload.",
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
	flags.StringVarP(&o.output, "output", "o", "name", "Output format. One of: (name, json, yaml)")
	flags.StringVar(&o.minAvailStr, "min-avail", "", "Minimum amount of available pods. Will be 1 if nothing else was set.")
	flags.StringVar(&o.maxUnavailStr, "max-unavail", "", "Maximum amount of unavailable Pods.")
	flags.BoolVar(&o.dryRun, "dry-run", false, "If set to true, it will only print the object, without sending it.")

	return cmd
}
