package engine

import (
	"testing"
)

func TestCopyCheckExecutor_Compute(t *testing.T) {
	executor := copyCheckExecutor{}
	param := Parameter{
		Path: []string{"./test_data"},
	}
	summary := executor.Compute(param, Config{CopyCheck: CopyCheckConfig{Threshold: 20}})

	if summary.Value != float64(1) {
		t.Error("run test copy check failed, the summary should equal be 1")
	}
}

func TestCopyCheckExecutorWithRegex_Compute(t *testing.T) {
	executor := copyCheckExecutor{}
	param := Parameter{
		Path: []string{"./test_data"},
	}
	summary := executor.Compute(param, Config{CopyCheck: CopyCheckConfig{Threshold: 20, IgnoreRegx: "copy_check"}})

	if summary.Value != float64(0) {
		t.Error("run test copy check failed, the summary should equal be 0")
	}
}
