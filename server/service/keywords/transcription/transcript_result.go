package transcription

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
)

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type TranscriptResult interface {
	ToStandard() StandardTranscriptResult
	GetText() string
}

type StandardTranscriptResult struct {
	Duration float64   `json:"duration"`
	Text     string    `json:"text"`
	Segments []Segment `json:"segments"`
}

func (s StandardTranscriptResult) ToJson() string {
	text, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println("error marshalling to json", err)
	}
	return string(text)
}

type AzureTranscriptResult struct {
	Duration                  string `json:"duration"`
	CombinedRecognizedPhrases []struct {
		Display string `json:"display"`
	} `json:"combinedRecognizedPhrases"`
	RecognizedPhrases []struct {
		Offset   string `json:"offset"`
		Duration string `json:"duration"`
		NBest    []struct {
			Display string `json:"display"`
		} `json:"nBest"`
	}
}

func (a AzureTranscriptResult) ToStandard() StandardTranscriptResult {
	return ConvertAzureToStandard(a)
}

func (a AzureTranscriptResult) GetText() string {
	return a.CombinedRecognizedPhrases[0].Display
}

func convertAzureDurationToSeconds(azureDuration string) float64 {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?`)
	matches := re.FindStringSubmatch(azureDuration)

	if matches == nil {
		log.Println("error parsing duration", azureDuration)
	}

	var totalSeconds float64

	if matches[1] != "" {
		hours, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Println("error parsing hours", err)
		}
		totalSeconds += float64(hours * 3600)
	}

	if matches[2] != "" {
		minutes, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Println("error parsing minutes", err)
		}
		totalSeconds += float64(minutes * 60)
	}

	if matches[3] != "" {
		seconds, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			log.Println("error parsing seconds", err)
		}
		totalSeconds += seconds
	}

	return totalSeconds
}

func ConvertAzureToStandard(azureResult AzureTranscriptResult) StandardTranscriptResult {
	var result StandardTranscriptResult
	result.Duration = convertAzureDurationToSeconds(azureResult.Duration)
	result.Text = azureResult.CombinedRecognizedPhrases[0].Display

	for _, phrase := range azureResult.RecognizedPhrases {
		stt := convertAzureDurationToSeconds(phrase.Offset)
		dur := convertAzureDurationToSeconds(phrase.Duration)
		seg := &Segment{
			Start: stt,
			End:   stt + dur,
			Text:  phrase.NBest[0].Display,
		}
		result.Segments = append(result.Segments, *seg)
	}
	return result
}
