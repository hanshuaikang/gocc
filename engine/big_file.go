package engine

import (
	"bufio"
	"os"
)

type bigFileExecutor struct {
}

func (e bigFileExecutor) computeFileLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	// nolint:errcheck
	defer f.Close()
	var lines int
	// 使用bufio.Scanner读取文件
	scanner := bufio.NewScanner(f)

	// 计数器

	// 逐行扫描文件
	for scanner.Scan() {
		lines++
	}
	return lines, nil
}

func (e bigFileExecutor) computeFilesLines(path string, config Config) (int, map[string]interface{}, error) {

	isDir, err := isDirectory(path)
	if err != nil {
		return 0, nil, err
	}
	bigFileNum := 0
	details := map[string]interface{}{}

	if isDir {
		goFiles, err := findGoFiles(path)
		if err != nil {
			return 0, nil, err
		}
		for _, file := range goFiles {
			lines, err := e.computeFileLines(file)
			if err != nil {
				return 0, nil, err
			}
			if lines > config.LintersSettings.BigFile.MaxLines {
				bigFileNum += 1
			}
			details[file] = lines
		}

		return bigFileNum, details, nil
	}
	lines, err := e.computeFileLines(path)
	if err != nil {
		return 0, nil, err
	}

	if lines > config.LintersSettings.BigFile.MaxLines {
		bigFileNum += 1
	}
	details[path] = lines
	return bigFileNum, details, nil

}

func (e bigFileExecutor) Compute(param Parameter, config Config) Summary {

	totalBigFileNum := 0
	var detailList []map[string]interface{}
	for _, path := range param.Path {
		bigFileNum, detail, err := e.computeFilesLines(path, config)
		if err != nil {
			return Summary{Name: BigFile, Err: err}
		}
		totalBigFileNum += bigFileNum
		detailList = append(detailList, detail)
	}

	details := mergeDetails(detailList)

	return Summary{Name: BigFile, Value: float64(totalBigFileNum), Details: details}

}
