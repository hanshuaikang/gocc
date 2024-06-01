package engine

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type unitTestExecutor struct {
}

func (e unitTestExecutor) parseOutPut(output string) (float64, error) {
	re := regexp.MustCompile(`coverage: (\d+\.\d+)% of statements`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		percentageStr := match[1]
		percentage, err := strconv.ParseFloat(percentageStr, 64)
		if err != nil {
			return 0, err
		}
		return percentage, nil
	}
	return 0, fmt.Errorf("no coverage information found in the output, output: %s", output)

}

func (e unitTestExecutor) covertPath(paths []string) ([]string, error) {

	var newPaths []string

	for _, path := range paths {
		isDir, err := isDirectory(path)
		if err != nil {
			return nil, err
		}
		if isDir {
			newPaths = append(newPaths, filepath.Join(path, "/..."))
			continue
		}
		newPaths = append(newPaths, path)
	}

	return newPaths, nil
}

// nolint
func (e unitTestExecutor) Compute(param Parameter, config Config) Summary {
	var out bytes.Buffer

	paths, err := e.covertPath(param.Path)
	if err != nil {
		return Summary{Name: UintTest, Err: err}
	}
	cmd := exec.Command("go", "test", strings.Join(paths, " "), "-cover")
	cmd.Stdout = &out
	err = cmd.Run()
	outPut := out.String()
	if err != nil {
		return Summary{Name: UintTest, Err: err}
	}

	if len(outPut) == 0 {
		return Summary{Name: UintTest, Err: fmt.Errorf("go test cmd output is empty")}
	}

	coverage, err := e.parseOutPut(outPut)
	if err != nil {
		return Summary{Name: UintTest, Err: err}
	}

	return Summary{
		Name:  UintTest,
		Value: coverage,
	}
}
