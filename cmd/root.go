package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gocc",
	Short: "gocc is a tool for Analyze the code quality of your GO project",
	Long:  `Analyze the code quality of your GO project, including the complexity of the circle, the grammatical errors ...`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func Execute() {
	_ = rootCmd.Execute()
}
