package engine

import (
	"path/filepath"
	"testing"
)

func TestUnitTest(t *testing.T) {

	unitExecutor := unitTestExecutor{}
	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/unit_test.go")
	param := Parameter{
		Path: []string{testFilePath},
	}
	summary := unitExecutor.Compute(param, Config{})

	if summary.Value != float64(0) {
		t.Errorf("unnitest run failed, expected value is 0, but value is %f", summary.Value)
	}
}
