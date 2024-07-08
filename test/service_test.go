package test

import (
	"bytes"
	"cognitube.com/keywords-service/server"
	"github.com/joho/godotenv"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestService(t *testing.T) {
	godotenv.Load("../dev.env")
	// Step 1: Open the local file
	localFile, err := os.Open("1 min.mp3")
	if err != nil {
		t.Fatalf("Failed to open local file: %v", err)
	}
	defer localFile.Close()

	var buf bytes.Buffer
	// Step 2: Create a new multipart writer
	multipartWriter := multipart.NewWriter(&buf)

	// Step 3: Create a form file part
	filePart, err := multipartWriter.CreateFormFile("audio", "file.mp3")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	// Step 4: Write the content of the local file to the form file part
	if _, err := io.Copy(filePart, localFile); err != nil {
		t.Fatalf("Failed to write local file to form file: %v", err)
	}

	// Step 5: Close the multipart writer
	if err := multipartWriter.Close(); err != nil {
		t.Fatalf("Failed to close multipart writer: %v", err)
	}

	// Step 6: Create a new http.Request
	r := httptest.NewRequest(http.MethodPost, "/api/getKeywords", &buf)
	r.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Extract the file part from the request
	file, _, err := r.FormFile("audio")
	if err != nil {
		t.Fatalf("Failed to extract file from request: %v", err)
	}

	// Step 8: Pass the multipart.File to GetKeyDescFromHttpAudioFile
	service := server.NewCognitubeKeywordsService()
	keywords, err := service.GetKeyDescFromHttpAudioFile(file)
	if err != nil {
		t.Fatalf("Error getting keywords: %v", err)
	}

	t.Logf("Keywords: %v", keywords)

}
