package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"cognitube.com/keywords-service/server/service/keywords"
)

type KeywordsHandler struct {
	keywordService keywords.ICognitubeKeywordsService
}

type RequestData struct {
	URL string `json:"url"`
}

func NewKeywordsHandler(service keywords.ICognitubeKeywordsService) *KeywordsHandler {
	return &KeywordsHandler{keywordService: service}
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
	/ - Home
	/api/v1/get-keywords - [GET] Get keywords
	/api/v1/transcription/create - [POST] Create transcription
	`))
}

func (h *KeywordsHandler) CallbackPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback POST request received")
	log.Printf("Request Url: %s", r.URL.String())
	log.Printf("Request Headers: %s", r.Header)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	log.Println(string(body))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(body) == 0 {
		// Register callback
		token := r.URL.Query().Get("validationToken")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
		return
	}

	go h.keywordService.OnTranscriptionCallback(body)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("received"))
}

func (h *KeywordsHandler) CreateTranscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Create transcription request received")
	var reqData struct {
		AudioURL string `json:"audioUrl"`
		VideoID  string `json:"videoId"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body, &reqData)
	log.Printf("Received request: %+v\n", reqData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := h.keywordService.CreateAsyncTranscription(reqData.AudioURL, reqData.VideoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("Created transcription with ID: %s\n", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (h *KeywordsHandler) ProcessTranscriptionResult(w http.ResponseWriter, r *http.Request) {
	log.Println("Process transcription result request received")
	var reqData struct {
		ID string `json:"id"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(body, &reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("Received request: %+v\n", reqData)
	h.keywordService.ProcessTranscriptionResult(reqData.ID)
	w.WriteHeader(http.StatusOK)
}

func (h *KeywordsHandler) GetKeywords(w http.ResponseWriter, r *http.Request) {
	log.Println("Get keywords request received")
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

	log.Printf("Downloaded and saved file: %+v\n", tempFile.Name())

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "Error seeking the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	keywords, err := h.keywordService.GetKeyDescFromHttpAudioFile(tempFile)
	if err != nil {
		http.Error(w, "Error getting keywords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Done processing the audio file to keywords")
	w.Write([]byte(keywords))
}
