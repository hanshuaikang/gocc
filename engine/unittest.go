package engine

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

func (e unitTestExecutor) computeCoverage() (float64, error) {

	var out bytes.Buffer

	cmd := exec.Command("go", "test", "./...", "-cover")
	cmd.Stdout = &out

	err := cmd.Run()
	outPut := out.String()
	if err != nil {
		return 0, errors.New(out.String())
	}

	if len(outPut) == 0 {
		return 0, fmt.Errorf("go test cmd output is empty")
	}

	coverage, err := e.parseOutPut(outPut)
	if err != nil {
		return 0, nil
	}

	return coverage, nil

}

// nolint
func (e unitTestExecutor) Compute(param Parameter, config Config) Summary {

	wd, err := getPwd()
	if err != nil {
		return Summary{Name: UintTest, Err: err}
	}

	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			panic(err)
		}
	}()

	details := map[string]interface{}{}

	totalCoverage := float64(0)

	for _, path := range param.Path {
		isGoProj, err := isGoProject(path)
		if err != nil {
			return Summary{Name: UintTest, Err: err}
		}

		if !isGoProj {
			continue
		}
		err = os.Chdir(path)
		if err != nil {
			return Summary{Name: UintTest, Err: err}
		}
		coverage, err := e.computeCoverage()
		if err != nil {
			return Summary{Name: UintTest, Err: err}
		}
		details[path] = coverage
		totalCoverage += coverage
	}

	var value float64
	if len(details) > 0 {
		value = totalCoverage / float64(len(details))
	}

	return Summary{
		Name:    UintTest,
		Value:   value,
		Details: details,
	}
}
