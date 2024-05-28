package engine

import (
	"bytes"
	"fmt"
	"os/exec"
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

// nolint
func (e unitTestExecutor) Compute(param Parameter, config Config) Summary {
	var out bytes.Buffer
	paths := strings.Join(param.Path, " ")
	cmd := exec.Command("go", "test", paths, "-cover")
	cmd.Stdout = &out
	err := cmd.Run()
	outPut := out.String()
	if err != nil {
		return Summary{Name: UintTest, Err: fmt.Errorf("err: %s, output: %s", err.Error(), outPut)}
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
