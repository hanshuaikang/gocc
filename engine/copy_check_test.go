package engine

import (
	"path/filepath"
	"testing"
)

func TestCopyCheckExecutor_Compute(t *testing.T) {
	executor := copyCheckExecutor{}

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/")

	param := Parameter{
		Path: []string{testFilePath},
	}
	summary := executor.Compute(param, Config{LintersSettings: LintersSettingsConfig{CopyCheck: CopyCheckConfig{Threshold: 20}}})

	if summary.Value != float64(1) {
		t.Error("run test copy check failed, the summary should equal be 1")
	}
}

func TestCopyCheckExecutorWithRegex_Compute(t *testing.T) {
	executor := copyCheckExecutor{}

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/")

	param := Parameter{
		Path: []string{testFilePath},
	}
	config := Config{LintersSettings: LintersSettingsConfig{CopyCheck: CopyCheckConfig{Threshold: 20, IgnoreRegx: "copy_check"}}}
	summary := executor.Compute(param, config)

	if summary.Value != float64(0) {
		t.Error("run test copy check failed, the summary should equal be 0")
	}
}
