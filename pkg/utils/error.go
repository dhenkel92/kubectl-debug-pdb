package utils

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

func HandleRunError(streams genericclioptions.IOStreams, err error) error {
	w := printers.GetNewTabWriter(streams.ErrOut)
	w.Write([]byte(err.Error()))
	return w.Flush()
}
