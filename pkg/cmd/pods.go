package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/dhenkel92/kubectl-utils/pkg/kube"
	printerUtils "github.com/dhenkel92/kubectl-utils/pkg/printer"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/rest"
)

// TODO: add sort-by argument
type PodsOptions struct {
	genericclioptions.IOStreams

	configFlags *genericclioptions.ConfigFlags
	Namespace   string
	PDBName     string
	Output      string
	Template    string

	OutputWide bool
}

func NewPodsOptions(streams genericclioptions.IOStreams, conf *genericclioptions.ConfigFlags) *PodsOptions {
	return &PodsOptions{
		IOStreams:   streams,
		configFlags: conf,
	}
}

func NewCmdPods(streams genericclioptions.IOStreams, conf *genericclioptions.ConfigFlags) *cobra.Command {
	o := NewPodsOptions(streams, conf)

	cmd := &cobra.Command{
		Use:   "pods [pdb name]",
		Short: "List pods",
		Long:  "List pods",
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
	flags.StringVarP(&o.Output, "output", "o", "human", "[human, json, yaml, jsonpath]")

	return cmd
}

func (o *PodsOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) < 1 {
		return fmt.Errorf("pdb name is required")
	}
	o.PDBName = args[0]

	o.Namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	if strings.HasPrefix(o.Output, "jsonpath") {
		split := strings.Split(o.Output, "=")
		o.Output = split[0]
		o.Template = split[1]
	}
	if o.Output == "wide" {
		o.OutputWide = true
	}
	return nil
}

func (o *PodsOptions) Validate() error {
	switch o.Output {
	case "json", "yaml", "human", "jsonpath", "wide":
	default:
		return fmt.Errorf("invalid output '%s'", o.Output)
	}

	return nil
}

func (o *PodsOptions) GetPrinter() (printers.ResourcePrinter, error) {
	switch o.Output {
	case "jsonpath":
		return printers.NewJSONPathPrinter(o.Template)
	case "json":
		return &printers.JSONPrinter{}, nil
	case "yaml":
		return &printers.YAMLPrinter{}, nil
	}

	tablePrinter := printers.NewTablePrinter(printers.PrintOptions{
		WithKind:      false,
		NoHeaders:     false,
		Wide:          o.OutputWide,
		WithNamespace: false,
		ShowLabels:    false,
		ColumnLabels:  []string{},
	})
	return &printerUtils.TablePrinter{Delegate: tablePrinter}, nil
}

func (o *PodsOptions) transformRequests(req *rest.Request) {
	if o.Output != "human" {
		return
	}
	req.SetHeader("Accept", strings.Join([]string{
		fmt.Sprintf("application/json;as=Table;v=%s;g=%s", metav1.SchemeGroupVersion.Version, metav1.GroupName),
		fmt.Sprintf("application/json;as=Table;v=%s;g=%s", metav1beta1.SchemeGroupVersion.Version, metav1beta1.GroupName),
		"application/json",
	}, ","))

	// TODO: sorting
	// if sorting, ensure we receive the full object in order to introspect its fields via jsonpath
	// if len(o.SortBy) > 0 {
	// req.Param("includeObject", "Object")
	// }
}

func (o *PodsOptions) Run() error {
	ctx := context.Background()
	printer, err := o.GetPrinter()
	if err != nil {
		return err
	}

	kubeClient, err := kube.New(o.configFlags)
	if err != nil {
		return err
	}

	pdb, err := kubeClient.GetClientset().PolicyV1().PodDisruptionBudgets(o.Namespace).Get(ctx, o.PDBName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println(err.Error())
			return nil
		}
		return err
	}

	selector, err := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
	if err != nil {
		return err
	}

	r := kubeClient.NewBuilder().Unstructured().
		NamespaceParam(o.Namespace).DefaultNamespace().
		ResourceTypeOrNameArgs(true, "pods").
		FieldSelectorParam("").
		LabelSelector(selector.String()).
		ContinueOnError().
		Latest().
		Flatten().
		TransformRequests(o.transformRequests).
		Do()

	if err := r.Err(); err != nil {
		return err
	}

	infos, err := r.Infos()
	if err != nil {
		return err
	}

	w := printers.GetNewTabWriter(o.Out)
	for _, pod := range infos {
		if err := printer.PrintObj(pod.Object, w); err != nil {
			return err
		}
	}

	return w.Flush()
}
