package test

import (
	"cognitube.com/keywords-service/transcription"
	"fmt"
	"os"
	"testing"
)

func TestWhisperClient_Transcript(t *testing.T) {
	// Hardcoding the environment variables for testing
	os.Setenv("APPLICATION_DEBUG", "true")
	os.Setenv("APPLICATION_WHISPER_URL", "https://api.openai.com/v1/audio/transcriptions")
	os.Setenv("APPLICATION_OPENAI_KEY", "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")

	client := transcription.NewWhisperClient()
	if client == nil {
		t.Errorf("Failed to create WhisperClient")
	}

	// Open the test audio file
	testFile, err := os.Open("./1 min.oga")
	if err != nil {
		t.Errorf("Failed to open test audio file")
	}

	res, err := client.Transcript(testFile)
	if err != nil || res == "" {
		t.Errorf("Failed to transcribe audio: %v", err)
	}
	fmt.Println(res)
}
