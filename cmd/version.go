package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version subcommand show git version info.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gocc V0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
