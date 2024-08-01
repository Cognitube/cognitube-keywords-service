package test

import (
	"cognitube.com/keywords-service/azure"
	"cognitube.com/keywords-service/transcription"
	"log"
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

func TestAzureClient_GetTranscriptionFileText(t *testing.T) {
	fileUrl := "https://spsvcprodeus.blob.core.windows.net/bestor-c6e3ae79-1b48-41bf-92ff-940bea3e5c2d/TranscriptionData/04e8e911-fe79-439a-92df-23510a39beef_0_0.json?skoid=50c6251a-ac54-47a3-9265-a1e4f84be9b9&sktid=33e01921-4d64-4f8c-a055-5bdaffd5e33d&skt=2024-07-19T06%3A05%3A55Z&ske=2024-07-24T06%3A10%3A55Z&sks=b&skv=2024-05-04&sv=2024-05-04&st=2024-07-19T06%3A05%3A55Z&se=2024-07-19T18%3A10%3A55Z&sr=b&sp=rl&sig=cr5Z3Ub3%2BQ2xunNPbvyW08aiqILhce5JMQbJEuFIEUU%3D"
	text, err := a.GetTranscriptionFileText(fileUrl)
	if err != nil {
		t.Error(err)
	}
	if text == "" {
		t.Error("Expected a non-empty text, got empty")
	}
	log.Println(text)
}
