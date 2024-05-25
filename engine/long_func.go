package engine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type longFuncExecutor struct {
}

func (e longFuncExecutor) computeFuncLength(path string, config Config) (int, map[string]interface{}, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return 0, nil, err
	}

	longFuncNum := 0
	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, path, content, parser.ParseComments)
	if err != nil {
		return 0, nil, err
	}

	details := map[string]interface{}{}

	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 计算函数体的行数
			bodyStart := fSet.Position(funcDecl.Body.Lbrace).Line
			bodyEnd := fSet.Position(funcDecl.Body.Rbrace).Line
			bodyLength := bodyEnd - bodyStart + 1
			// example: engine/big_file.go Compute 72:89
			key := fmt.Sprintf("%s %s %d:%d", path, funcDecl.Name, bodyStart, bodyEnd)
			if bodyLength > config.LongFunc.MaxLength {
				longFuncNum += 1
			}
			details[key] = bodyLength
		}
	}
	return longFuncNum, details, nil
}

func (e longFuncExecutor) computeFuncLengths(path string, config Config) (int, map[string]interface{}, error) {

	isD, err := isDir(path)
	if err != nil {
		return 0, nil, nil
	}
	longFuncNum := 0

	if isD {
		goFiles, err := findGoFiles(path)
		var detailsList []map[string]interface{}
		if err != nil {
			return 0, nil, nil
		}
		for _, file := range goFiles {
			num, d, err := e.computeFuncLength(file, config)
			if err != nil {
				return 0, nil, nil
			}
			longFuncNum += num
			detailsList = append(detailsList, d)
		}
		return longFuncNum, mergeDetails(detailsList), nil
	}

	// single file
	num, details, err := e.computeFuncLength(path, config)
	if err != nil {
		return 0, nil, nil
	}
	return num, details, nil

}

func (e longFuncExecutor) Compute(param Parameter, config Config) Summary {

	var detailList []map[string]interface{}
	longFuncNum := 0
	for _, path := range param.Path {
		num, detail, err := e.computeFuncLengths(path, config)
		if err != nil {
			return Summary{Name: LongFunc, ErrMsg: err.Error()}
		}
		longFuncNum += num
		detailList = append(detailList, detail)
	}

	return Summary{Name: LongFunc, Details: mergeDetails(detailList), Value: float64(longFuncNum)}

}
