package server

import (
	"fmt"
	"net/http"
)

type KeywordsHandler struct {
	Service ICognitubeKeywordsService
}

func NewKeywordsHandler(service ICognitubeKeywordsService) *KeywordsHandler {
	return &KeywordsHandler{Service: service}
}

func (h *KeywordsHandler) GetKeywords(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Error retrieving the audio file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Printf("Received File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	keywords, err := h.Service.GetKeyDescFromHttpAudioFile(file)
	if err != nil {
		http.Error(w, "Error getting keywords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Done processing the audio file ", handler.Filename, " to keywords")
	w.Write([]byte(keywords))
}
