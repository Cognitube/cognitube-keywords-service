package transcription

import (
	"os"
)

type Transcriber interface {
	// Transcript an audio from audio file to text
	Transcript(file *os.File) (string, error)
}

func NewTranscriber(name string) Transcriber {
	if name == "whisper" {
		return NewWhisperClient()
	}
	return nil
}
