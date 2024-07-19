package transcription

import (
	"cognitube.com/keywords-service/env"
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

type IAzureClient interface {
	GetJobIDFromSelfURL(url string) string
	GetAllTranscriptionFileURLs(jobId string) ([]string, error)
	GetTranscriptionFileText(fileUrl string) (string, error)
}

type AzureClient struct {
	IAzureClient
	ApiKey        string
	BatchTransURL string
}

func (a *AzureClient) OnTranscriptionCallback(r *http.Request) (*TranscriptionResult, error) {
	// Get current finished job ID
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	url := gjson.GetBytes(body, "self").String()
	id := a.GetJobIDFromSelfURL(url)

	// list all transcription files and get the first one's URL
	listFileResp, err := a.GetAllTranscriptionFileURLs(id)
	if err != nil {
		return nil, err
	}
	transcriptionJsonURL := listFileResp[0]

	// GET the file and extract the text
	text, err := a.GetTranscriptionFileText(transcriptionJsonURL)
	if err != nil {
		return nil, err
	}

	// Build the result
	result := TranscriptionResult{
		ID:   id,
		Text: text,
	}
	fmt.Println("Transcription result: ", result.ID, result.Text[:50])
	return &result, nil
}

func (a *AzureClient) GetJobIDFromSelfURL(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func (a *AzureClient) GetAllTranscriptionFileURLs(jobId string) ([]string, error) {
	var buffer string
	err := requests.
		URL(a.BatchTransURL+"/"+jobId+"/files").
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", a.ApiKey).
		ToString(&buffer).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	arr := gjson.Get(buffer, "values.#.links.contentUrl").Array()
	res := make([]string, len(arr))
	for i := range arr {
		res[i] = arr[i].String()
	}
	return res, nil
}

func (a *AzureClient) GetTranscriptionFileText(fileUrl string) (string, error) {
	var buffer string
	err := requests.
		URL(fileUrl).
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", a.ApiKey).
		ToString(&buffer).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	return gjson.Get(buffer, "combinedRecognizedPhrases.0.display").String(), nil
}

func (a *AzureClient) CreateTranscription(fileUrl string, displayName string) (string, error) {
	var buffer string

	err := requests.
		URL(a.BatchTransURL).
		Method(http.MethodPost).
		ContentType("application/json").
		Header("Ocp-Apim-Subscription-Key", a.ApiKey).
		BodyJSON(map[string]interface{}{
			"contentUrls": []string{fileUrl},
			"locale":      "en-US",
			"displayName": displayName,
			"model":       nil,
			"properties": map[string]interface{}{
				"diarizationEnabled":                    false,
				"wordLevelTimestampsEnabled":            false,
				"displayFormWordLevelTimestampsEnabled": false,
				"punctuationMode":                       "DictatedAndAutomatic",
				"profanityFilterMode":                   "Masked",
			}}).
		ToString(&buffer).
		Fetch(context.Background())

	if err != nil {
		return "", err
	}

	url := gjson.Get(buffer, "self").String()
	return a.GetJobIDFromSelfURL(url), nil
}

func NewAzureClient() AsyncTranscriber {
	return &AzureClient{
		ApiKey:        env.GetInstance().AzureKey,
		BatchTransURL: env.GetInstance().AzureUrl,
	}
}
