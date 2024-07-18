package transcription

import (
	"net/http"
)

type TranscriptionResult struct {
	ID      string
	Text    string
	FileUrl string
}

type AsyncTranscriber interface {
	CreateTranscription(fileUrl string, displayName string) (string, error)
	OnTranscriptionCallback(r *http.Request) (*TranscriptionResult, error)
}

func NewAsyncTranscriberClient(name string) AsyncTranscriber {
	if name == "azure" {
		return NewAzureClient()
	}
	return nil
}
