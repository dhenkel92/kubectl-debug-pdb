package cmd

import (
	"context"
	"fmt"

	"github.com/dhenkel92/kubectl-pdb/pkg/kube"
	"github.com/spf13/cobra"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

type EvictPodOptions struct {
	genericclioptions.IOStreams

	configFlags *genericclioptions.ConfigFlags
	Namespace   string
	PodName     string
	DryRun      bool
	Output      string
}

func NewEvictPodOptions(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *EvictPodOptions {
	return &EvictPodOptions{
		IOStreams:   streams,
		configFlags: configFlags,
	}
}

func NewCmdEvictPod(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	o := NewEvictPodOptions(streams, configFlags)

	cmd := &cobra.Command{
		Use:   "pod",
		Short: "Utility to evict a pod from a node",
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

	cmd.Flags().StringVarP(&o.Output, "output", "o", "json", "Output format. One of: json|yaml")
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", true, "If true, only print the object that would be sent, without sending it.")

	return cmd
}

func (o *EvictPodOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) != 1 {
		return fmt.Errorf("pod name is required")
	}
	o.PodName = args[0]

	o.Namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	return nil
}

func (o *EvictPodOptions) GetPrinter() (printers.ResourcePrinter, error) {
	switch o.Output {
	case "json":
		return &printers.JSONPrinter{}, nil
	case "yaml":
		return &printers.YAMLPrinter{}, nil
	}

	return nil, fmt.Errorf("invalid output '%s'", o.Output)
}

func (o *EvictPodOptions) Validate() error {
	switch o.Output {
	case "json", "yaml":
	default:
		return fmt.Errorf("invalid output '%s'", o.Output)
	}

	return nil
}

func (o *EvictPodOptions) Run() error {
	ctx := context.Background()

	kubeClient, err := kube.New(o.configFlags)
	if err != nil {
		return err
	}

	pod, err := kubeClient.GetClientset().CoreV1().Pods(o.Namespace).Get(ctx, o.PodName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	dryRunOpts := []string{}
	if o.DryRun {
		dryRunOpts = append(dryRunOpts, "All")
	}

	eviction := policyv1.Eviction{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Eviction",
			APIVersion: policyv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: &metav1.DeleteOptions{
			DryRun: dryRunOpts,
		},
	}

	err = kubeClient.GetClientset().CoreV1().Pods(o.Namespace).EvictV1(ctx, &eviction)
	if err != nil {
		return err
	}

	w := printers.GetNewTabWriter(o.Out)
	printer, err := o.GetPrinter()
	if err := printer.PrintObj(&eviction, w); err != nil {
		return err
	}

	return w.Flush()
}
