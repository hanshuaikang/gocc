package engine

import (
	"fmt"
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
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func matchRegex(str, pattern string) (bool, error) {
	if len(pattern) == 0 {
		return false, nil
	}
	// Compile the regular expression pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %v", err)
	}
	// Check if the filename matches the regex pattern
	matches := re.MatchString(str)
	return matches, nil
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
	return regexp.Compile(expr)
}

// round retain the latter two of the float64 value
func round(num float64) float64 {
	return math.Round(num*100) / 100
}
