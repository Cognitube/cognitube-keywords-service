package transcription

type AsyncTranscriber interface {
	CreateTranscription(fileUrl string, displayName string) (string, error)
	OnTranscriptionCallback(id string) (string, error)
}

func NewAsyncTranscriberClient(name string) AsyncTranscriber {
	if name == "azure" {
		return NewAzureClient()
	}
	return nil
}
