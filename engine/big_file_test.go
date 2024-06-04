package engine

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBigFileExecutor_Compute(t *testing.T) {

	executor := bigFileExecutor{}

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/")
	param := Parameter{
		Path: []string{testFilePath},
	}
	config := Config{LintersSettings: LintersSettingsConfig{BigFile: BigFileConfig{MaxLines: 800}}}

	summary := executor.Compute(param, config)

	if summary.Value != 0 {
		t.Errorf(" bigfile executor compute failed, the summary should equal 0")
	}
}

func TestBigFileExecutor_Compute_Single_File(t *testing.T) {

	executor := bigFileExecutor{}
	baseDir, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}
	testFilePath := filepath.Join(baseDir, "test_data/bigfile.go")
	param := Parameter{
		Path: []string{testFilePath},
	}
	config := Config{LintersSettings: LintersSettingsConfig{BigFile: BigFileConfig{MaxLines: 1}}}
	summary := executor.Compute(param, config)

	if int(summary.Value) != 1 {
		t.Errorf("bigfile executor compute failed, the summary value should equal 1")
	}

	if summary.Details[testFilePath].(int) != 8 {
		t.Errorf("bigfile executor compute failed, the test file's lines should equal 7")
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
