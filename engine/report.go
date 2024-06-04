package engine

import (
	"fmt"
	"strconv"
)

type Summary struct {
	Name     string                 `json:"name"`
	Value    float64                `json:"value"`
	Details  map[string]interface{} `json:"details"`
	Duration float64                `json:"duration"`
	Err      error                  `json:"error"`
}

func (s Summary) Print() {

	switch s.Name {
	case CyclomaticComplexity, BigFile, LongFunc:
		for k, v := range s.Details {
			fmt.Printf("%s value: %d\n", k, v.(int))
		}
	case UintTest:
		for k, v := range s.Details {
			fmt.Printf("%s Coverage: %.2f\n", k, v.(float64))
		}
	case CopyCheck:
		for k, v := range s.Details {
			index, _ := strconv.Atoi(k)
			fmt.Printf("Found Copy #%d\n", index+1)
			for _, i := range v.([]string) {
				println(i)
			}
		}
	case Security:
		for k, v := range s.Details {
			values, ok := v.([]string)
			if !ok {
				continue
			}
			fmt.Printf("Found %d Vulnerability in %s\n", len(values), k)
			for _, item := range values {
				fmt.Println(item)
			}

		}
	case Syntax:
		for k, v := range s.Details {
			fmt.Printf("%s %s\n", k, v.(string))
		}
	}

}
