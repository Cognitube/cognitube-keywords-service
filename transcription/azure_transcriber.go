package transcription

import (
	"cognitube.com/keywords-service/azure"
	"cognitube.com/keywords-service/env"
	"encoding/json"
)

type AzureTranscriber struct {
	azure.SpeechClient
}

func (a *AzureTranscriber) CreateTranscription(fileUrl string, displayName string) (string, error) {
	// POST the file to Azure and get the transcription ID
	return a.SpeechClient.CreateTranscription(fileUrl, displayName)
}

func (a *AzureTranscriber) OnTranscriptionCallback(id string) (string, error) {
	// list all transcription files and get the first one's URL
	listFileResp, err := a.GetAllTranscriptionFileURLs(id)
	if err != nil {
		return "", err
	}
	transcriptionJsonURL := listFileResp[0]

	// GET the file and extract the text
	content, err := a.GetTranscriptionFileContent(transcriptionJsonURL)
	if err != nil {
		return "", err
	}

	var result AzureTranscriptResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}

	standard := ConvertAzureToStandard(result)
	text, err := json.MarshalIndent(standard, "", "  ")
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func NewAzureClient() AsyncTranscriber {
	return &AzureTranscriber{
		SpeechClient: azure.SpeechClient{
			ApiKey:        env.GetInstance().AzureKey,
			BatchTransURL: env.GetInstance().AzureUrl,
		},
	}
}
