package test

import (
	"cognitube.com/keywords-service/azure"
	"cognitube.com/keywords-service/transcription"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"
)

var a = transcription.AzureTranscriber{
	SpeechClient: azure.SpeechClient{
		ApiKey:        "HIDE_THIS_BEFORE_COMMITTING_TO_GITHUB",
		BatchTransURL: "https://eastus.api.cognitive.microsoft.com/speechtotext/v3.1/transcriptions",
	},
}

func TestAzureClient_GetJobIDFromSelfURL(t *testing.T) {
	url := "https://eastus.api.cognitive.microsoft.com/speechtotext/v3.1/transcriptions/04e8e911-fe79-439a-92df-23510a39beef"
	id := azure.GetJobIDFromSelfURL(url)
	expected := "04e8e911-fe79-439a-92df-23510a39beef"
	if id != expected {
		t.Errorf("Expected %s, got %s", expected, id)
	}

}

func TestAzureClient_GetAllTranscriptionFileURLs(t *testing.T) {
	jobId := "04e8e911-fe79-439a-92df-23510a39beef"
	urls, err := a.GetAllTranscriptionFileURLs(jobId)
	if err != nil {
		t.Error(err)
	}
	if len(urls) == 0 {
		t.Error("Expected at least one URL, got none")
	}
	for _, url := range urls {
		if url == "" {
			t.Error("Expected a non-empty URL, got empty")
		}
		log.Println(url)
	}
}

func TestAdapter(t *testing.T) {
	f, err := os.Open("azure_transcript.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(string(content))
	var result transcription.AzureTranscriptResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(result)

	standardResult := transcription.ConvertAzureToStandard(result)
	//t.Log(standardResult)
	jsText, err := json.MarshalIndent(standardResult, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsText))
}
