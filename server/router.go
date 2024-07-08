package server

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the HTTP routes for the server.
func SetupRoutes(r *mux.Router, keywordsService ICognitubeKeywordsService) {
	keywordsHandler := NewKeywordsHandler(keywordsService)
	r.HandleFunc("/keywords", keywordsHandler.GetKeywords).Methods("POST")
}
