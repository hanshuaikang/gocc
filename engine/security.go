package engine

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/google/go-cmdtest"
	"golang.org/x/vuln/scan"
)

type securityExecutor struct {
}

func (e securityExecutor) parseVulnerabilityOutput(output string) []string {
	outputLines := strings.Split(output, "\n")

	var vulnerabilities []string
	var vulnerabilityItem []string
	matched := false
	for _, line := range outputLines {

		if strings.HasPrefix(line, "Vulnerability ") {
			matched = true
		}

		if len(line) == 0 {
			matched = false
			if len(vulnerabilityItem) != 0 {
				vulnerabilities = append(vulnerabilities, strings.Join(vulnerabilityItem, "\n"))
				vulnerabilityItem = []string{}

			}
		}

		if matched {
			vulnerabilityItem = append(vulnerabilityItem, line)
		}
	}

	return vulnerabilities
}

func (e securityExecutor) scanGoProject(path string, envs []string) (int, map[string]interface{}, error) {

	buf := &bytes.Buffer{}
	cmd := scan.Command(context.Background(), []string{"-C", path, "./..."}...)
	cmd.Stdout = buf
	cmd.Stderr = buf

	// config
	cmd.Env = append(os.Environ(), envs...)
	if err := cmd.Start(); err != nil {
		return 0, nil, err
	}
	err := cmd.Wait()
	switch e := err.(type) {
	case nil:
	case interface{ ExitCode() int }:
		err = &cmdtest.ExitCodeErr{Msg: err.Error(), Code: e.ExitCode()}
		if e.ExitCode() == 0 || e.ExitCode() == 3 {
			err = nil
		}

	default:
		err = &cmdtest.ExitCodeErr{Msg: err.Error(), Code: 1}

	}

	if err != nil {
		return 0, nil, err
	}

	vulnerabilities := e.parseVulnerabilityOutput(buf.String())
	details := map[string]interface{}{}
	details[path] = vulnerabilities
	return len(vulnerabilities), details, nil
}

func (e securityExecutor) Compute(param Parameter, config Config) Summary {

	vulnerabilityNum := 0
	var detailList []map[string]interface{}
	for _, path := range param.Path {
		isDir, err := isDirectory(path)
		if err != nil {
			return Summary{Name: Security, Err: err}
		}
		if !isDir {
			continue
		}
		isGoPro, err := isGoProject(path)
		if err != nil {
			return Summary{Name: Security, Err: err}
		}

		if !isGoPro {
			continue
		}
		num, detail, err := e.scanGoProject(path, config.LintersSettings.Security.Env)
		if err != nil {
			return Summary{Name: Security, Err: err}
		}

		vulnerabilityNum += num
		detailList = append(detailList, detail)
	}

	return Summary{Name: Security, Value: float64(vulnerabilityNum), Details: mergeDetails(detailList)}
}
