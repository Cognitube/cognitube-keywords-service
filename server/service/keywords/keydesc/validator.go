package keydesc

import (
	"encoding/json"
)

type KeywordsItem struct {
	Keyword     string  `json:"keyword"`
	Description string  `json:"description"`
	StartTime   float64 `json:"startTime"`
	EndTime     float64 `json:"endTime"`
}

func ValidateKeywordsResult(jsonStr string) bool {
	var s []KeywordsItem
	err := json.Unmarshal([]byte(jsonStr), &s)
	if err != nil {
		return false
	}
	return true
}
