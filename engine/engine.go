package engine

type Executor interface {
	Compute(param Parameter, config Config) Summary
}

func RunAllTools(param Parameter, config Config) []Summary {
	var summaryList []Summary
	for _, executor := range executeHub {
		summary := executor.Compute(param, config)
		summaryList = append(summaryList, summary)
	}

	return summaryList

}
