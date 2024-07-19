package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type KeywordsHandler struct {
	Service ICognitubeKeywordsService
}

type RequestData struct {
	URL string `json:"url"`
}

func NewKeywordsHandler(service ICognitubeKeywordsService) *KeywordsHandler {
	return &KeywordsHandler{Service: service}
}

func (h *KeywordsHandler) CallbackGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Callback GET request received")
	token := r.URL.Query().Get("validationToken")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (h *KeywordsHandler) CallbackPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Callback POST request received")
	go h.Service.OnTranscriptionCallback(r)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("received"))
}

func (h *KeywordsHandler) CreateTranscription(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create transcription request received")
	var reqData struct {
		AudioURL    string `json:"audio_url"`
		DisplayName string `json:"display_name"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body, &reqData)
	fmt.Printf("Received request: %+v\n", reqData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := h.Service.CreateAsyncTranscription(reqData.AudioURL, reqData.DisplayName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("Created transcription with ID: %s\n", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (h *KeywordsHandler) GetKeywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get keywords request received")
	var requestData RequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error decoding JSON from body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.URL == "" {
		http.Error(w, "URL parameter is missing in the JSON body", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(requestData.URL)
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
