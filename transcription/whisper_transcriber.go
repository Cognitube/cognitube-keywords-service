package transcription

import (
	"bytes"
	"cognitube.com/keywords-service/env"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type WhisperTranscriber struct {
	apiKey     string
	whisperUrl string
}

func NewWhisperClient() *WhisperTranscriber {
	return &WhisperTranscriber{
		apiKey:     env.GetInstance().OpenAIKey,
		whisperUrl: env.GetInstance().WhisperUrl,
	}
}

type WhisperResponse struct {
	start float64
	end   float64
	text  string
}

func (c *WhisperTranscriber) Transcript(file *os.File) (string, error) {
	if !env.GetInstance().Debug {
		defer os.Remove(file.Name())
	}
	defer file.Close()

	// Create a buffer and a multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file part
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file contents: %w", err)
	}

	// Add other form fields
	_ = writer.WriteField("model", "whisper-1")
	_ = writer.WriteField("response_format", "verbose_json")

	// Close the multipart writer to set the terminating boundary
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.whisperUrl, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	apiKey := env.GetInstance().OpenAIKey

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read and decode the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-ok status from API: %d; response: %s", resp.StatusCode, responseData)
	}

	// Extract "segments" field as a JSON string
	var responseMap map[string]json.RawMessage
	if err := json.Unmarshal(responseData, &responseMap); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	rawMsg, exists := responseMap["text"]
	if !exists {
		return "", fmt.Errorf("segments field not found")
	}

	//var transcriptObj map[string]interface{}
	//if err := json.Unmarshal(responseData, &transcriptObj); err != nil {
	//	return "", fmt.Errorf("error parsing transcript: %w", err)
	//}
	//
	//extractedTranscript, ok := transcriptObj["segments"].(string)
	//if !ok {
	//	if errorMsg, ok := transcriptObj["error"].(map[string]interface{}); ok {
	//		return "", fmt.Errorf("error in transcript: %s", errorMsg["message"])
	//	}
	//	return "", fmt.Errorf("segments field not found or invalid")
	//}

	return string(rawMsg), nil
}
