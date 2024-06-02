package engine

import (
	"path/filepath"
	"testing"
)

func TestSecurityExecutor_Compute(t *testing.T) {
	executor := securityExecutor{}

	current, err := getPwd()
	if err != nil {
		t.Error(err.Error())
	}

	summary := executor.Compute(Parameter{Path: []string{filepath.Dir(current)}}, Config{})

	if summary.Value <= float64(0) {
		t.Errorf("run security failed， the value obtained should be greater than 0")
	}

}
