package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hanshuaikang/gocc/engine"
	"gopkg.in/yaml.v2"
)

const usageDoc = `Analyze the code quality of your GO project, including the complexity of the circle, the grammatical errors ...
Usage:
	gocc [flags] <Go file or directory> ...
Flags:
	-config  the config is a yaml file, you can provide configuration files to meet some personalized configurations

the output is a json file, example: https://github.com/hanshuaikang/gocc
`

func main() {

	configFile := flag.String("config", "", "the gocc config file path")

	log.SetFlags(0)
	log.SetPrefix("gocc: ")
	flag.Usage = usage

	flag.Parse()
	paths := flag.Args()
	var err error
	config := engine.DefaultConfig()
	if *configFile != "" {
		config, err = parseConfig(*configFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "gcc run err : load config yaml failed: %s", err.Error())
			os.Exit(1)
		}
	}

	if len(paths) == 0 {
		usage()
	}

	param := engine.Parameter{Path: paths}
	summaryList := engine.RunAllTools(param, config)
	printOutputToFile(summaryList)

}

func printOutputToFile(summaryList []engine.Summary) {
	output, err := json.MarshalIndent(summaryList, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "gcc run err : print output to json file failed, json marshal err : %s", err.Error())
		os.Exit(1)

	}
	err = writeFile(output)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "gcc run err : print output to json file failed, wirte json file err : %s", err.Error())
		os.Exit(1)
	}
}

func parseConfig(path string) (engine.Config, error) {
	// 读取 YAML 文件
	file, err := os.Open(path)
	if err != nil {
		return engine.Config{}, err
	}
	defer file.Close()

	var config engine.Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return engine.Config{}, err
	}
	return config, nil
}

func writeFile(output []byte) error {

	file, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(output))
	if err != nil {
		return err
	}

	return err

}

func usage() {
	_, _ = fmt.Fprint(os.Stderr, usageDoc)
	os.Exit(2)
}
