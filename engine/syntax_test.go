package engine

import (
	"path/filepath"
	"testing"
)

func TestSyntaxExecutor_Compute(t *testing.T) {

	baseDir, err := getPwd()
	if err != nil {
		t.Error(err)
	}

	executor := syntaxExecutor{}
	param := Parameter{Path: []string{filepath.Dir(baseDir)}}

	summary := executor.Compute(param, Config{})

	if summary.Value != float64(0) {
		t.Error("run syntax failed, the value should equal o")
	}

}
