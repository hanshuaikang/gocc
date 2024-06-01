package engine

import (
	"path/filepath"
	"testing"
)

func TestLongFuncExecutor_Compute(t *testing.T) {
	executor := longFuncExecutor{}

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}

	param := Parameter{Path: []string{
		filepath.Join(baseDir, "test_data/big_func.go"),
		filepath.Join(baseDir, "test_data/"),
	}}
	config := Config{LintersSettings: LintersSettingsConfig{LongFunc: LongFuncConfig{MaxLength: 1}}}
	summary := executor.Compute(param, config)
	if summary.Value != float64(9) {
		t.Errorf("big file executor execute failed, the summary value show equal %f != 9.0", summary.Value)
	}
}
