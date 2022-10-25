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

	rootCmd := &cobra.Command{
		Use:   "pdb",
		Short: "Utility to work with pod disruption budgets",
	}
	rootCmd.AddCommand(cmd.NewCmdPdb(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
