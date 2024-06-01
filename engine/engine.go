package engine

import (
	"time"
)

type Executor interface {
	Compute(param Parameter, config Config) Summary
}

func isEnabled(enabled []string, executorName string) bool {
	for _, v := range enabled {
		if v == executorName {
			return true
		}
	}
	return false
}

func RunAllTools(param Parameter, config Config) []Summary {
	var summaryList []Summary
	for name, executor := range executeHub {
		if len(config.Linters.Enable) > 0 && !isEnabled(config.Linters.Enable, name) {
			continue
		}
		start := time.Now()
		summary := executor.Compute(param, config)
		summary.Duration = time.Since(start).Seconds()
		summaryList = append(summaryList, summary)
	}

	return summaryList

}
