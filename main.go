package main

import (
	"cognitube.com/keywords-service/server"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Use -env flag to specify the .env file to load
	envFile := flag.String("env", ".env", "The name of the .env file to load")
	flag.Parse()

	err := godotenv.Load(*envFile)
	if err != nil {
		fmt.Println("No .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	keywordsServer := server.NewCongitubeKeywordsServer()
	keywordsServer.StartListening(port)
}
