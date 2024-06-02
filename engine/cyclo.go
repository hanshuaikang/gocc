package engine

import (
	"fmt"
	"regexp"

	"github.com/fzipp/gocyclo"
)

type cyclomaticComplexityExecutor struct {
}

func (c cyclomaticComplexityExecutor) buildDetails(stats gocyclo.Stats, config Config) map[string]interface{} {

	details := map[string]interface{}{}

	for _, stat := range stats {
		if stat.Complexity < config.LintersSettings.Cyclo.Over {
			continue
		}
		// example: engine (bigFileExecutor).Compute engine/big_file.go:72:1
		key := fmt.Sprintf("%s %s %s", stat.PkgName, stat.FuncName, stat.Pos)
		details[key] = stat.Complexity
	}
	return details

}

func (c cyclomaticComplexityExecutor) Compute(param Parameter, config Config) Summary {

	var re *regexp.Regexp
	var err error
	if len(config.LintersSettings.Cyclo.IgnoreRegx) != 0 {
		re, err = regex(config.LintersSettings.Cyclo.IgnoreRegx)
		if err != nil {
			return Summary{Name: CyclomaticComplexity, Err: err}
		}
	}

	stats := gocyclo.Analyze(param.Path, re)

	details := c.buildDetails(stats, config)
	summary := Summary{
		Name:    CyclomaticComplexity,
		Value:   round(stats.AverageComplexity()),
		Details: details,
	}

	return summary

}
