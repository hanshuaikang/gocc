package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hanshuaikang/gocc/engine"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run subcommand start analyze the code",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := engine.DefaultConfig()
		configFilePath := cmd.Flag("config").Value.String()
		if len(configFilePath) != 0 {
			config, err = parseConfig(configFilePath)
			if err != nil {
				Error(cmd, args, err)
			}
		}
		if len(args) == 0 {
			Error(cmd, args, fmt.Errorf("path is empty"))
		}

		param := engine.Parameter{Path: covertAbsPath(args)}
		summaryList := engine.RunAllTools(param, config)
		report(summaryList, config.ReportType)
	},
}

func covertAbsPath(paths []string) []string {
	var absPaths []string
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		if err != nil {
			Error(nil, paths, fmt.Errorf("path can't not covert to abs path"))
		}
		absPaths = append(absPaths, abs)
	}
	return absPaths
}

func parseConfig(path string) (engine.Config, error) {
	// 读取 YAML 文件
	file, err := os.Open(path)
	// nolint:errcheck
	if err != nil {
		return engine.Config{}, err
	}
	// nolint:errcheck
	defer file.Close()

	var config engine.Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return engine.Config{}, err
	}
	return config, nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
