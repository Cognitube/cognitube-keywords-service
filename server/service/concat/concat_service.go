package concat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ConcatJSON struct {
	Keyword     string `json:"keyword"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Clip        int    `json:"clip"`
}

type ConcatResponse struct {
	Results []ConcatJSON
}

type IConcatService interface {
	ConcatJSONByOrder(reqs []ConcatJSON) (ConcatResponse, error)
}

type ConcatService struct {
}

func NewConcatService() IConcatService {
	return &ConcatService{}
}

func (s *ConcatService) ConcatJSONByOrder(reqs []ConcatJSON) (ConcatResponse, error) {
	sort.Slice(reqs, func(i, j int) bool {
		clipI := reqs[i].Clip
		clipJ := reqs[j].Clip
		if clipI != clipJ {
			return clipI < clipJ
		}
		timeI, _ := parseTime(reqs[i].Start)
		timeJ, _ := parseTime(reqs[j].Start)
		return timeI < timeJ
	})

	return ConcatResponse{Results: reqs}, nil
}

func parseTime(t string) (time.Duration, error) {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
