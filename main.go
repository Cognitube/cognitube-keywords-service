package main

import (
	"flag"
	"fmt"
	"os"

	"cognitube.com/keywords-service/server"
	"github.com/joho/godotenv"
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
		port = "8082"
	}

	keywordsServer := server.NewCongitubeKeywordsServer()
	keywordsServer.StartListening(port)
}
