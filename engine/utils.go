package engine

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
)

// fondGoFiles found all go files of a path
func findGoFiles(dirPath string) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" {
			goFiles = append(goFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return goFiles, nil
}

// isDir determine whether a path is a folder
func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// mergeDetails merge a list of details
func mergeDetails(detailList []map[string]interface{}) map[string]interface{} {
	details := map[string]interface{}{}

	for _, detail := range detailList {
		for k, v := range detail {
			details[k] = v
		}
	}
	return details
}

// regex str cover to regexp.Regexp object
func regex(expr string) (*regexp.Regexp, error) {
	if expr == "" {
		return nil, nil
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return re, nil
}

// round retain the latter two of the float64 value
func round(num float64) float64 {
	return math.Round(num*100) / 100
}
