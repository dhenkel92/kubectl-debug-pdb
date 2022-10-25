package cmd

import (
	"fmt"
	"strings"

	"github.com/dhenkel92/kubectl-pdb/pkg/kube"
	"github.com/liggitt/tabwriter"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

type PDBOptions struct {
	genericclioptions.IOStreams

	configFlags   *genericclioptions.ConfigFlags
	AllNamespaces bool
	Namespace     string
	PodName       string
	Output        string
	Template      string
}

func NewPDBOptions(streams genericclioptions.IOStreams) *PDBOptions {
	return &PDBOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func (o *PDBOptions) GetPrinter() (printers.ResourcePrinter, error) {
	switch o.Output {
	case "jsonpath":
		return printers.NewJSONPathPrinter(o.Template)
	case "json":
		return &printers.JSONPrinter{}, nil
	case "yaml":
		return &printers.YAMLPrinter{}, nil
	}
	return printers.NewTablePrinter(printers.PrintOptions{
		WithKind:      false,
		NoHeaders:     false,
		Wide:          false,
		WithNamespace: o.AllNamespaces,
		ShowLabels:    false,
		ColumnLabels:  []string{},
	}), nil
}

func (o *PDBOptions) GetNamespace() string {
	if o.AllNamespaces {
		return ""
	}
	return o.Namespace
}

func (o *PDBOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	o.Namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	if len(args) > 0 && args[0] != "" {
		o.PodName = args[0]
	}

	if strings.HasPrefix(o.Output, "jsonpath") {
		split := strings.Split(o.Output, "=")
		o.Output = split[0]
		o.Template = split[1]
	}
	return nil
}

func (o *PDBOptions) Validate() error {
	switch o.Output {
	case "json", "yaml", "human", "jsonpath":
	default:
		return fmt.Errorf("invalid output '%s'", o.Output)
	}
	return nil
}

func (o *PDBOptions) getWriter() *tabwriter.Writer {
	return printers.GetNewTabWriter(o.Out)
}

func (o *PDBOptions) Run() error {
	printer, err := o.GetPrinter()
	if err != nil {
		return err
	}

	kubeClients, err := kube.New(o.configFlags)
	if err != nil {
		return err
	}

	pods, err := kubeClients.GetNamespacedPods(o.GetNamespace(), o.PodName)
	if err != nil {
		return err
	}

	pdbs, err := kubeClients.GetNamespacedPDBs(o.GetNamespace())
	if err != nil {
		return err
	}

	entries := make([]PodPDBEntry, 0, len(pods))
	for ns, podLst := range pods {
		podPDBs, ok := pdbs[ns]
		if !ok {
			continue
		}

		for _, pod := range podLst {
			pdbs := getMatchingPDBs(&pod, podPDBs)
			entries = append(entries, PodPDBEntry{
				Pod:  pod,
				Pdbs: pdbs,
			})
		}
	}

	lst := &PodPDBList{Items: entries}
	w := o.getWriter()
	if err := printer.PrintObj(lst.toMetaTable(), w); err != nil {
		return err
	}
	return w.Flush()
}

func NewCmdPdb(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewPDBOptions(streams)

	cmd := &cobra.Command{
		Use:   "cover [pod name]",
		Short: "Shows which PDBs are covering the given workload.",
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
	flags.BoolVarP(&o.AllNamespaces, "all-namespaces", "a", false, "")
	flags.StringVarP(&o.Output, "output", "o", "human", "[human, json, yaml, jsonpath]")
	o.configFlags.AddFlags(flags)

	return cmd
}

func getMatchingPDBs(pod *corev1.Pod, pdbs []policyv1.PodDisruptionBudget) []policyv1.PodDisruptionBudget {
	ls := labels.Set(pod.Labels)

	matchingPDBs := make([]policyv1.PodDisruptionBudget, 0)
	for _, pdb := range pdbs {
		// TODO catch error
		selector, _ := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
		if selector.Matches(ls) {
			matchingPDBs = append(matchingPDBs, pdb)
		}
	}
	return matchingPDBs
}
