package server

import (
	"mime/multipart"
	"net/http"
)

// ICognitubeKeywordsService receives the direct input from the HTTP request and returns the output to the HTTP response.
// Considering HTTP response time. Do not perform time-consuming operations outside the scope of the interface
// Move the time-consuming operations to the implementation of the interface, and if the return value is not necessary for the HTTP response, use goroutines
type ICognitubeKeywordsService interface {
	GetKeyDescFromHttpAudioFile(file multipart.File) (string, error)
	CreateAsyncTranscription(fileUrl string, displayName string) (string, error)
	OnTranscriptionCallback(r *http.Request)
}
