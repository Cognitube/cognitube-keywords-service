package keywords_service

import (
	"cognitube.com/keywords-service/server/service/keywords/transcription"
	"os"
	"testing"
)

func TestWhisper_Transcribe(t *testing.T) {
	os.Setenv("APPLICATION_DEBUG", "true")
	os.Setenv("APPLICATION_OPENAI_KEY", "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")
	os.Setenv("APPLICATION_WHISPER_URL", "https://api.openai.com/v1/audio/transcriptions")
	trans := transcription.NewTranscriber("whisper")
	filename := "../static/test.oga"
	// open file
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	text, err := trans.Transcript(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(text)
}
