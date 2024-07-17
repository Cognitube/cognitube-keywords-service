package main

import (
	"os"

	"cognitube.com/keywords-service/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	keywordsServer := server.NewCongitubeKeywordsServer()
	keywordsServer.StartListening(port)
}
