package handler

import (
	"cognitube.com/keywords-service/server/service/concat"
	"encoding/json"
	"net/http"
)

type ConcatHandler struct {
	concatService concat.IConcatService
}

func (h *ConcatHandler) ConcatJSONByOrder(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var keywords []concat.ConcatJSON
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			var fileKeywords concat.ConcatJSON
			if err := json.NewDecoder(file).Decode(&fileKeywords); err != nil {
				http.Error(w, "Error decoding JSON from file: "+err.Error(), http.StatusBadRequest)
				return
			}
			keywords = append(keywords, fileKeywords)
		}
	}
	var res, _ = h.concatService.ConcatJSONByOrder(keywords)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res.Results); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewConcatHandler(concatService concat.IConcatService) *ConcatHandler {
	return &ConcatHandler{
		concatService: concatService,
	}
}
