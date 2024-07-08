package server

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the HTTP routes for the server.
func SetupRoutes(r *mux.Router, keywordsService ICognitubeKeywordsService) {
	keywordsHandler := NewKeywordsHandler(keywordsService)
	r.HandleFunc("/api/v1/get-keywords", keywordsHandler.GetKeywords).Methods("POST")
}
