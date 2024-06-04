package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hanshuaikang/gocc/engine"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run subcommand start analyze the code",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		configFilePath := cmd.Flag("config").Value.String()
		config, err := getConfig(configFilePath)
		if err != nil {
			Error(cmd, args, err)
		}

		if len(args) == 0 {
			Error(cmd, args, fmt.Errorf("path is empty"))
		}

		param := engine.Parameter{Path: covertAbsPath(args)}
		summaryList := engine.RunAllTools(param, config)
		report(summaryList, config.ReportType, config.IgnoreError)
	},
}

func covertAbsPath(paths []string) []string {
	var absPaths []string
	for _, path := range paths {
		file, err := os.Stat(path)
		if err != nil {
			Error(nil, paths, fmt.Errorf("path can't not covert to abs path"))
		}

		if !file.IsDir() && !strings.HasSuffix(path, ".go") {
			continue
		}

		abs, err := filepath.Abs(path)

		if err != nil {
			Error(nil, paths, fmt.Errorf("path can't not covert to abs path"))
		}
		absPaths = append(absPaths, abs)
	}
	return absPaths
}

func getConfig(path string) (engine.Config, error) {
	if len(path) == 0 {
		// 检查当前工作目录是否存在 gocc.yaml 配置文件

		wd, err := os.Getwd()
		if err != nil {
			return engine.Config{}, err
		}

		filePath := filepath.Join(wd, "gocc.yaml")

		_, err = os.Stat(filePath)
		if err != nil {
			return engine.DefaultConfig(), nil
		}

		return parseConfig(filePath)
	}

	return parseConfig(path)
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
