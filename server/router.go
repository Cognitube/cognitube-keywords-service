package server

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the HTTP routes for the server.
func SetupRoutes(r *mux.Router, keywordsService ICognitubeKeywordsService) {
	keywordsHandler := NewKeywordsHandler(keywordsService)
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/api/v1/get-keywords", keywordsHandler.GetKeywords).Methods("POST") // Sync
	r.HandleFunc("/api/v1/callback", keywordsHandler.CallbackPost).Methods("POST")
	r.HandleFunc("/api/v1/transcription/create", keywordsHandler.CreateTranscription).Methods("POST")
}
