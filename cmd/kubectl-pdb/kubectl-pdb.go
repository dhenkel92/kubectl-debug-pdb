package main

import (
	"os"

	"github.com/dhenkel92/kubectl-utils/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-pdb", pflag.ExitOnError)
	pflag.CommandLine = flags

	conf := genericclioptions.NewConfigFlags(true)
	conf.AddFlags(flags)

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	rootCmd := &cobra.Command{
		Use:   "pdb",
		Short: "Utility to work with pod disruption budgets",
	}
	rootCmd.AddCommand(cmd.NewCmdPdb(streams, conf))
	rootCmd.AddCommand(cmd.NewCmdCreatePDB(streams, conf))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
