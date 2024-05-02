package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const whisperUrl = "https://api.openai.com/v1/audio/transcriptions"

func extractAudio(videoFilePath string) (string, error) {
	// Determine the output file path
	dir := filepath.Dir(videoFilePath)
	audioFilePath := filepath.Join(dir, "extracted_audio.mp4")

	// Construct the ffmpeg command to extract audio
	cmd := exec.Command("ffmpeg", "-i", videoFilePath, "-vn", "-acodec", "aac", "-b:a", "96000", "-ac", "1", "-ar", "44100", audioFilePath)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %w", err)
	}

	return audioFilePath, nil
}

func getKeywordsFromGemini(transcript string) (string, error) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")

	model.SetTemperature(1)
	model.SetTopP(0.95)
	model.SetTopK(0)
	model.SetMaxOutputTokens(65536)

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("You are a keyword extraction model. \n Please extract 20 - 30 keywords from the following user context, where each keyword should be low frequency nouns that is important to the user context. The user input format is a list of json file, where the start and the text attributes in segment are very useful in this problem. your output should looks like this: \n[{keyword:\\\"keyword1\\\", description:\\\"description of keyword1 in the given context\\\",startTime:0.0,endTime:3.83}, {keyword:\\\"keyword2\\\", explain:\\\"description of keyword2 in the given context\\\",startTime:3.83,endTime:5.17}, ...]. The description of the keywords should be specific, first explain the actual meaning of this keyword, and then related to the context explain what this keyword means in the video, and no less than 150 words. The startTime and endTime of the return value should round to two decimal places.\n\n")},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(transcript))
	if err != nil {
		log.Fatal(err)
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")

	return "", nil
}

func getTranscriptFromWhisper(audioFilePath string) (string, error) {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", err
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
	req, err := http.NewRequest("POST", whisperUrl, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

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

	segmentsData, exists := responseMap["segments"]
	if !exists {
		return "", fmt.Errorf("segments field not found")
	}

	return string(segmentsData), nil
}

func getKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form containing the file
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the video file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file to save the uploaded video
	tempFile, err := os.CreateTemp("", "temp-*.mp4")
	if err != nil {
		http.Error(w, "Error creating a temporary file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	log.Printf("Temporary File: %+v\n", tempFile.Name())
	// Copy the file to the destination
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)

	log.Printf("File saved to: %s\n", tempFile.Name())
	// Convert video to audio
	audioPath, err := extractAudio(tempFile.Name())
	if err != nil {
		http.Error(w, "Error converting video to audio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Audio saved to: %s\n", audioPath)
	// Get transcript from Whisper API
	transcript, err := getTranscriptFromWhisper(audioPath)
	if err != nil {
		http.Error(w, "Error getting transcript from Whisper API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Transcript: %s\n", transcript)
	keywords, err := getKeywordsFromGemini(transcript)
	if err != nil {
		http.Error(w, "Error getting keywords from Gemini API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the transcript
	w.Write([]byte(keywords))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define a handler function
	http.HandleFunc("/getKeywords", getKeywordsHandler)

	// Start the server on localhost port 8080
	log.Println("Starting server on :8083")
	err = http.ListenAndServe(":8083", nil) // nil tells it to use the default router we set up with http.HandleFunc
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
