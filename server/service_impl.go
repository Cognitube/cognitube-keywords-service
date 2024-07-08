package server

import (
	"cognitube.com/keywords-service/keydesc"
	"cognitube.com/keywords-service/transcription"
	"io"
	"mime/multipart"
	"os"
)

type CongnitubeKeywordsService struct {
	descriptor  keydesc.KeywordsDescriptor
	transcriber transcription.Transcriber
}

func NewCognitubeKeywordsService() ICognitubeKeywordsService {
	return &CongnitubeKeywordsService{
		descriptor:  keydesc.NewKeywordsDescriptor("gpt"),
		transcriber: transcription.NewTranscriber("whisper"),
	}
}

func convertMultipartFileToOsFile(file multipart.File) (*os.File, error) {
	// Convert multipart.File to os.File
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "uploaded-*.mp3")
	if err != nil {
		return nil, err
	}

	// Copy the content of the multipart file to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close() // Attempt to close the file in case of error
		return nil, err
	}

	// It's important to seek to the beginning of the file before returning it
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close() // Attempt to close the file in case of error
		return nil, err
	}

	return tempFile, nil
}

func (c *CongnitubeKeywordsService) GetKeyDescFromHttpAudioFile(multiFile multipart.File) (string, error) {
	// Convert multipart.File to os.File
	defer multiFile.Close()
	file, err := convertMultipartFileToOsFile(multiFile)
	if err != nil {
		return "", err
	}

	// Transcribe the audio file to text
	text, err := c.transcriber.Transcript(file)
	if err != nil {
		return "", err
	}

	// Describe the text to keywords
	keywords, err := c.descriptor.Describe(text)
	if err != nil {
		return "", err
	}
	return keywords, nil
}
