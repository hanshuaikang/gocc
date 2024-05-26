package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Error(cmd *cobra.Command, args []string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}
