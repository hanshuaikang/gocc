package engine

type Summary struct {
	Name    string                 `json:"name"`
	Value   float64                `json:"value"`
	Details map[string]interface{} `json:"details"`
	Err     error                  `json:"error"`
}
