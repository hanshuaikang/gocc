package engine

import (
	"testing"
)

func TestUnitTest(t *testing.T) {

	unitExecutor := unitTestExecutor{}

	param := Parameter{
		Path: []string{"./test_data/unit_test.go"},
	}
	summary := unitExecutor.Compute(param, Config{})

	if summary.Value != float64(0) {
		t.Errorf("unnitest run failed, expected value is 0, but value is %f", summary.Value)
	}
}
