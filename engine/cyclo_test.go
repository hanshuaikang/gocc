package engine

import (
	"testing"
)

func TestCyclo(t *testing.T) {

	cyclo := cyclomaticComplexityExecutor{}

	param := Parameter{
		Path: []string{"./test_data/cyclo.go"},
	}
	summary := cyclo.Compute(param, Config{})

	if summary.Value != float64(2) {
		t.Errorf("cyclomatic complexity compute faild, the summary should equal 2")
	}
}
