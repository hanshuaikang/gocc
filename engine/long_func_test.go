package engine

import "testing"

func TestLongFuncExecutor_Compute(t *testing.T) {
	executor := longFuncExecutor{}

	param := Parameter{Path: []string{
		"./test_data/big_func.go",
		"./test_data/",
	}}

	summary := executor.Compute(param, Config{LongFunc: LongFuncConfig{MaxLength: 1}})
	if summary.Value != float64(8) {
		t.Error("big file executor execute failed, the summary value show equal 8")
	}
}
