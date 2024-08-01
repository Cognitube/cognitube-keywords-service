package main

import (
	"cognitube.com/keywords-service/azure"
	"cognitube.com/keywords-service/env"
	"os"

	"cognitube.com/keywords-service/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	azure.RegisterCallback(env.GetInstance().ApplicationCallbackUrl)
	keywordsServer := server.NewCongitubeKeywordsServer()
	keywordsServer.StartListening(port)
}
