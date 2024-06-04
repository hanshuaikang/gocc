package engine

import (
	"path/filepath"
	"testing"
)

func TestCyclo(t *testing.T) {

	cyclo := cyclomaticComplexityExecutor{}

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/cyclo.go")

	param := Parameter{
		Path: []string{testFilePath},
	}
	summary := cyclo.Compute(param, Config{})

	if summary.Value != float64(2) {
		t.Errorf("cyclomatic complexity compute faild, the summary should equal 2")
	}
}
