package transcription

import (
	"cognitube.com/keywords-service/env"
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
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
	apiKey   string
	azureUrl string
}

func (a *AzureClient) OnTranscriptionCallback(r *http.Request) (*TranscriptionResult, error) {
	// Get current finished job ID
	var resp JobJSON
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	id := a.GetJobIDFromSelfURL(resp.SelfUrl)

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

	// Notify the listener
	result := TranscriptionResult{
		ID:   id,
		Text: text,
	}
	fmt.Println("Transcription result: ", result.ID, result.Text[:50])
	return &result, nil
}

type JobJSON struct {
	SelfUrl string `json:"self"`
}

func (a *AzureClient) GetJobIDFromSelfURL(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func (a *AzureClient) GetAllTranscriptionFileURLs(jobId string) ([]string, error) {
	var listFileResp struct {
		Values []struct {
			Kind string `json:"kind"`
			Name string `json:"name"`
			Url  string `json:"self"`
		}
	}
	err := requests.
		URL(a.azureUrl+"/transcriptions/"+jobId+"/files").
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", a.apiKey).
		ToJSON(&listFileResp).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	var urls []string
	for _, file := range listFileResp.Values {
		urls = append(urls, file.Url)
	}
	return urls, nil
}

func (a *AzureClient) GetTranscriptionFileText(fileUrl string) (string, error) {
	var resultFileJSON struct {
		CombinedRecognizedPhrases []struct {
			Display string `json:"display"`
		}
	}
	err := requests.
		URL(fileUrl).
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", a.apiKey).
		ToJSON(&resultFileJSON).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	return resultFileJSON.CombinedRecognizedPhrases[0].Display, nil
}

func (a *AzureClient) CreateTranscription(fileUrl string, displayName string) (string, error) {

	var res JobJSON

	err := requests.
		URL(a.azureUrl).
		Method(http.MethodPost).
		ContentType("application/json").
		Header("Ocp-Apim-Subscription-Key", a.apiKey).
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
		ToJSON(&res).
		Fetch(context.Background())

	if err != nil {
		return "", err
	}

	return a.GetJobIDFromSelfURL(res.SelfUrl), nil
}

func NewAzureClient() AsyncTranscriber {
	return &AzureClient{
		apiKey:   env.GetInstance().AzureKey,
		azureUrl: env.GetInstance().AzureUrl,
	}
}
