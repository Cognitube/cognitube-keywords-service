package azure

import (
	"bytes"
	"context"
	"github.com/carlmjohnson/requests"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"strings"
)

func GetJobIDFromSelfURL(url string) string {
	if !strings.Contains(url, "/") {
		return ""
	}
	return url[strings.LastIndex(url, "/")+1:]
}

type ISpeechClient interface {
	GetAllTranscriptionFileURLs(jobId string) ([]string, error)
	GetTranscriptionFileContent(fileUrl string) ([]byte, error)
	RegisterCallback(callbackUrl string) error
}

type SpeechClient struct {
	ISpeechClient
	BatchTransURL string
	ApiKey        string
}

func (c *SpeechClient) GetAllTranscriptionFileURLs(jobId string) ([]string, error) {
	var buffer string
	err := requests.
		URL(c.BatchTransURL+"/"+jobId+"/files").
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", c.ApiKey).
		ToString(&buffer).
		Fetch(context.Background())
	if err != nil {
		log.Println("Failed to get transcription files: ", err.Error())
		return nil, err
	}

	arr := gjson.Get(buffer, "values.#.links.contentUrl").Array()
	res := make([]string, len(arr))
	for i := range arr {
		res[i] = arr[i].String()
	}
	return res, nil
}

func (c *SpeechClient) GetTranscriptionFileContent(fileUrl string) ([]byte, error) {
	var buffer bytes.Buffer
	err := requests.
		URL(fileUrl).
		Method(http.MethodGet).
		Header("Ocp-Apim-Subscription-Key", c.ApiKey).
		ToBytesBuffer(&buffer).
		Fetch(context.Background())
	if err != nil {
		log.Println("Failed to get transcription file text: ", err.Error())
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *SpeechClient) CreateTranscription(fileUrl string, displayName string) (string, error) {
	var buffer string

	err := requests.
		URL(c.BatchTransURL).
		Method(http.MethodPost).
		ContentType("application/json").
		Header("Ocp-Apim-Subscription-Key", c.ApiKey).
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
		log.Println("Failed to create transcription: ", err.Error())
		return "", err
	}

	url := gjson.Get(buffer, "self").String()
	return GetJobIDFromSelfURL(url), nil
}
