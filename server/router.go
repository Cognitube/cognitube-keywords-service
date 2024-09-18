package server

import (
	"cognitube.com/keywords-service/server/handler"
	"github.com/gorilla/mux"
)

// SetupRoutes configures the HTTP routes for the server.
func SetupRoutes(r *mux.Router, s *CongitubeAIServer) {
	// Keywords Service Routes
	keywordsHandler := handler.NewKeywordsHandler(s.keywordsService)
	r.HandleFunc("/", handler.Home).Methods("GET")
	r.HandleFunc("/api/v1/get-keywords", keywordsHandler.GetKeywords).Methods("POST") // Sync
	r.HandleFunc("/api/v1/callback", keywordsHandler.CallbackPost).Methods("POST")    // This callback is broken
	r.HandleFunc("/api/v1/transcription/create", keywordsHandler.CreateTranscription).Methods("POST")
	r.HandleFunc("/api/v1/transcription/process", keywordsHandler.ProcessTranscriptionResult).Methods("POST")

	// Concatenation Service Routes
	concatHandler := handler.NewConcatHandler(s.concatService)
	r.HandleFunc("/api/v1/concat", concatHandler.ConcatJSONByOrder).Methods("POST")
}
