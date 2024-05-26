package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hanshuaikang/gocc/engine"
)

func report(summaryList []engine.Summary, reportType string) {
	switch reportType {
	case engine.Json:
		reportJsonFile(summaryList)
	case engine.Console:
		reportConsole(summaryList)
	default:
		fmt.Fprintf(os.Stderr, "gcc run err : unsupport report type : %s", reportType)
	}
}

func reportJsonFile(summaryList []engine.Summary) {
	output, err := json.MarshalIndent(summaryList, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "gcc run err : print output to json file failed, json marshal err : %s", err)
		os.Exit(1)

	}
	err = writeFile(output, "output.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "gcc run err : print output to json file failed, wirte json file err : %s", err)
		os.Exit(1)
	}
}

func reportConsole(summaryList []engine.Summary) {
	fmt.Printf("code analysis is complete. Below is the brief output. Please switch to HTML or JSON output type for detailed output.\n\n")
	for _, summary := range summaryList {
		if summary.Err != nil {
			fmt.Printf("run %s failed, err: %v\n\n", summary.Name, summary.Err)
			continue
		}
		fmt.Printf("Name: %s\n", summary.Name)
		fmt.Printf("Value: %v\n\n", summary.Value)
	}

}

func writeFile(output []byte, fileName string) error {
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(output))

	return err

}
