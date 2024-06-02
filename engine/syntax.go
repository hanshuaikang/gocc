package engine

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type syntaxExecutor struct {
}

func (e syntaxExecutor) parseOutPut(path string, output string) (int, map[string]interface{}) {

	detail := map[string]interface{}{}
	count := 0
	// 定义正则表达式模式，匹配 .go:<数字>:<数字>:
	pattern := regexp.MustCompile(`\S+\.go:\d+:\d+:.*`)

	// 查找所有匹配的子字符串
	matches := pattern.FindAllString(output, -1)

	// 输出结果
	for _, match := range matches {
		count += 1
		key := filepath.Join(path, strings.Split(match, " ")[0])
		detail[key] = strings.Join(strings.Split(match, " ")[1:], " ")

	}
	return count, detail
}

func (e syntaxExecutor) runVet(path string) (int, map[string]interface{}) {

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "vet", "./...")
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	// go vet 命令如果 err, 则说明有语法错误
	count, detail := e.parseOutPut(path, errOut.String())
	if err != nil {
		return count, detail
	}

	return 0, map[string]interface{}{}

}

func (e syntaxExecutor) Compute(param Parameter, config Config) Summary {

	wd, err := getPwd()
	if err != nil {
		return Summary{Name: Syntax, Err: err}
	}

	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			panic(err)
		}
	}()

	syntaxCount := 0
	var details []map[string]interface{}
	for _, path := range param.Path {
		isGoPro, err := isGoProject(path)
		if err != nil {
			return Summary{Name: Syntax, Err: err}
		}

		if !isGoPro {
			continue
		}

		err = os.Chdir(path)
		if err != nil {
			return Summary{Name: Syntax, Err: err}
		}

		count, detail := e.runVet(path)
		syntaxCount += count
		details = append(details, detail)
	}

	return Summary{
		Name: Syntax, Details: mergeDetails(details), Value: float64(syntaxCount),
	}
}
