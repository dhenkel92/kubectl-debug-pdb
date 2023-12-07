package main

import (
	"os"

	"github.com/dhenkel92/kubectl-debug-pdb/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	// https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-evict", pflag.ExitOnError)
	pflag.CommandLine = flags

	conf := genericclioptions.NewConfigFlags(true)
	conf.AddFlags(flags)

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	rootCmd := &cobra.Command{
		Use:   "evict",
		Short: "Utility to evict a pod from a node",
	}
	rootCmd.AddCommand(cmd.NewCmdEvictPod(streams, conf))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
