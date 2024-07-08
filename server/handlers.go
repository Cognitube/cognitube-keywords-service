package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type KeywordsHandler struct {
	Service ICognitubeKeywordsService
}

func NewKeywordsHandler(service ICognitubeKeywordsService) *KeywordsHandler {
	return &KeywordsHandler{Service: service}
}

func (h *KeywordsHandler) GetKeywords(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error downloading file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Non-OK HTTP status: "+resp.Status, http.StatusInternalServerError)
		return
	}

	tempFile, err := os.CreateTemp("", "download-*.oga")
	if err != nil {
		http.Error(w, "Error creating a temp file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		http.Error(w, "Error saving the downloaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Downloaded and saved file: %+v\n", tempFile.Name())

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "Error seeking the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	keywords, err := h.Service.GetKeyDescFromHttpAudioFile(tempFile)
	if err != nil {
		http.Error(w, "Error getting keywords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Done processing the audio file to keywords")
	w.Write([]byte(keywords))
}
