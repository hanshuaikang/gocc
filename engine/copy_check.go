package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mibk/dupl/job"
	"github.com/mibk/dupl/syntax"
)

type Clone struct {
	Filename  string `json:"filename"`
	LineStart int    `json:"lineStart"`
	LineEnd   int    `json:"lineEnd"`
	Fragment  []byte `json:"fragment"`
}

// group by file name
type byNameAndLine []Clone

func (c byNameAndLine) Len() int { return len(c) }

func (c byNameAndLine) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c byNameAndLine) Less(i, j int) bool {
	if c[i].Filename == c[j].Filename {
		return c[i].LineStart < c[j].LineStart
	}
	return c[i].Filename < c[j].Filename
}

type copyCheckExecutor struct{}

// crawlPaths get all go files of a path
func (e copyCheckExecutor) crawlPaths(paths []string, config Config) chan string {
	fChan := make(chan string)
	go func() {
		defer close(fChan)
		for _, path := range paths {
			matched, err := matchRegex(path, config.LintersSettings.CopyCheck.IgnoreRegx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "copy check run failed, match regex error, please check your regex")
				break
			}
			if matched {
				continue
			}
			info, err := os.Lstat(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "copy check run failed: crawl path err, skip this file :%s, err: %s", path, err)
				continue
			}
			if !info.IsDir() {
				fChan <- path
				continue
			}
			err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				// the IgnoreRegx is checked, not throw err
				matched, _ = matchRegex(path, config.LintersSettings.CopyCheck.IgnoreRegx)

				if matched {
					return nil
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
					fChan <- path
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "copy check run failed: crawl path err, skip this dir :%s, err: %s", path, err)
				continue
			}
		}
	}()
	return fChan
}

// copyCheckExecutor Calculate the start and end of the code block
func (e copyCheckExecutor) blockLines(file []byte, from, to int) (int, int, error) {

	if from < 0 || to > len(file) || from >= to {
		return 0, 0, fmt.Errorf("compute block lines failed, the from|to is invalid, from: %d, to: %d", from, to)
	}

	line := 1
	lineStart, lineEnd := 0, 0
	for offset, b := range file {
		if b == '\n' {
			line++
		}
		if offset == from {
			lineStart = line
		}
		if offset == to-1 {
			lineEnd = line
			break
		}
	}

	// Ensure lineEnd is set in case the loop completes without breaking
	if lineEnd == 0 {
		lineEnd = line
	}

	return lineStart, lineEnd, nil
}

func (e copyCheckExecutor) prepareClonesInfo(duplicateNodes [][]*syntax.Node) ([]Clone, error) {
	clones := make([]Clone, len(duplicateNodes))
	for i, dup := range duplicateNodes {
		cnt := len(dup)
		if cnt == 0 {
			continue
		}
		start := dup[0]
		end := dup[cnt-1]

		file, err := os.ReadFile(start.Filename)
		if err != nil {
			return nil, err
		}

		cl := Clone{Filename: start.Filename}
		cl.LineStart, cl.LineEnd, err = e.blockLines(file, start.Pos, end.End)
		if err != nil {
			return nil, err
		}
		clones[i] = cl
	}
	return clones, nil
}

func (e copyCheckExecutor) buildDetails(dupChan <-chan syntax.Match) (map[string]interface{}, error) {

	groups := make(map[string][][]*syntax.Node)
	for dup := range dupChan {
		groups[dup.Hash] = append(groups[dup.Hash], dup.Frags...)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	details := map[string]interface{}{}

	for i, k := range keys {
		uniq := e.unique(groups[k])
		if len(uniq) <= 1 {
			continue
		}
		clones, err := e.prepareClonesInfo(uniq)
		if err != nil {
			return nil, err
		}
		sort.Sort(byNameAndLine(clones))
		details[strconv.Itoa(i)] = clones
	}

	return details, nil
}

// unique Remove duplicates
func (e copyCheckExecutor) unique(group [][]*syntax.Node) [][]*syntax.Node {
	fileMap := make(map[string]map[int]struct{})

	var newGroup [][]*syntax.Node
	for _, seq := range group {
		node := seq[0]
		file, ok := fileMap[node.Filename]
		if !ok {
			file = make(map[int]struct{})
			fileMap[node.Filename] = file
		}
		if _, ok := file[node.Pos]; !ok {
			file[node.Pos] = struct{}{}
			newGroup = append(newGroup, seq)
		}
	}
	return newGroup
}

func (e copyCheckExecutor) Compute(param Parameter, config Config) Summary {

	pathChan := job.Parse(e.crawlPaths(param.Path, config))
	// generate ast tree
	t, data, done := job.BuildTree(pathChan)
	<-done
	// finish stream
	t.Update(&syntax.Node{Type: -1})

	matchChan := t.FindDuplOver(config.LintersSettings.CopyCheck.Threshold)
	dupChan := make(chan syntax.Match)

	go func() {
		for m := range matchChan {
			match := syntax.FindSyntaxUnits(*data, m, config.LintersSettings.CopyCheck.Threshold)
			if len(match.Frags) > 0 {
				dupChan <- match
			}
		}
		close(dupChan)
	}()

	details, err := e.buildDetails(dupChan)
	if err != nil {
		return Summary{Name: CopyCheck, Err: err}
	}

	return Summary{Name: CopyCheck, Details: details, Value: float64(len(details))}

}
