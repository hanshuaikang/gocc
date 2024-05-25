package engine

import (
	"strings"
	"testing"
)

func TestBigFileExecutor_Compute(t *testing.T) {

	executor := bigFileExecutor{}

	param := Parameter{
		Path: []string{"./test_data/"},
	}
	summary := executor.Compute(param, Config{BigFile: BigFileConfig{MaxLines: 800}})

	if summary.Value != 0 {
		t.Errorf(" bigfile executor compute failed, the summary should equal 0")
	}

	if summary.Details["test_data/bigfile.go"].(int) != 7 {
		t.Errorf("bigfile executor compute failed, the test file's lines should equal 11")
	}
}

func TestBigFileExecutor_Compute_Single_File(t *testing.T) {

	executor := bigFileExecutor{}

	param := Parameter{
		Path: []string{"./test_data/bigfile.go"},
	}
	summary := executor.Compute(param, Config{BigFile: BigFileConfig{MaxLines: 1}})

	if int(summary.Value) != 1 {
		t.Errorf("bigfile executor compute failed, the summary value should equal 1")
	}

	if summary.Details["./test_data/bigfile.go"].(int) != 7 {
		t.Errorf("bigfile executor compute failed, the test file's lines should equal 11")
	}

}

func TestBigFileExecutor_Compute_Failed(t *testing.T) {

	executor := bigFileExecutor{}

	param := Parameter{
		Path: []string{"./xxxxx/"},
	}
	summary := executor.Compute(param, Config{})

	if summary.Err == nil {
		t.Errorf(" bigfile executor compute faild, the result should error")
	}

	if !strings.EqualFold(summary.Err.Error(), "stat ./xxxxx/: no such file or directory") {
		t.Errorf(" bigfile executor compute failed, the errMsg should equal 'stat ./xxxxx/: no such file or directory' ")
	}

}
