package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hanshuaikang/gocc/engine"
)

func report(summaryList []engine.Summary, reportType string, ignoreError bool) {
	switch reportType {
	case engine.Json:
		reportJsonFile(summaryList, ignoreError)
	case engine.Console:
		reportConsole(summaryList, ignoreError)

	default:
		fmt.Fprintf(os.Stderr, "gcc run err : unsupport report type : %s", reportType)
	}
}

// nolint
func reportJsonFile(summaryList []engine.Summary, ignoreError bool) {
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

func reportConsole(summaryList []engine.Summary, ignoreError bool) {
	fmt.Printf("code analysis is complete. Below is the brief output. Please switch to HTML or JSON output type for detailed output.\n\n")

	fmt.Print("Summary: ")
	for _, summary := range summaryList {
		fmt.Printf("%s: %.2f  ", summary.Name, summary.Value)
	}

	hasError := false
	for _, summary := range summaryList {
		fmt.Printf("\nName: %s\n", summary.Name)
		if summary.Err != nil {
			fmt.Printf("Err: run %s failed, err: %v\n", summary.Name, summary.Err)
			fmt.Printf("Duration: %vs\n\n", summary.Duration)
			hasError = true
			continue
		}

		// 这些指标如果存在 detail，则说明有不合格的地方
		switch summary.Name {
		case engine.Security, engine.BigFile, engine.Syntax, engine.CyclomaticComplexity, engine.LongFunc, engine.CopyCheck:
			if len(summary.Details) > 0 {
				hasError = true
			}
		}

		fmt.Printf("Value: %v\n", summary.Value)
		fmt.Printf("Duration: %3.fs\n", summary.Duration)
		summary.Print()
	}
	if hasError && !ignoreError {
		os.Exit(1)
	}
}

func writeFile(output []byte, fileName string) error {
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}
	// nolint:errcheck
	defer file.Close()

	_, err = io.WriteString(file, string(output))

	return err

}
