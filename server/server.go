package server

import (
	"log"
	"net/http"

	"cognitube.com/keywords-service/server/service/concat"
	"cognitube.com/keywords-service/server/service/keywords"

	"github.com/gorilla/mux"
)

type WebServer interface {
	StartListening(port string)
}

type CognitubeAIServer struct {
	keywordsService keywords.ICognitubeKeywordsService
	concatService   concat.IConcatService
}

func (c *CognitubeAIServer) StartListening(port string) {
	router := mux.NewRouter()

	SetupRoutes(router, c)

	log.Println("Server started at port ", port)

	http.ListenAndServe(":"+port, router)
}

func NewCognitubeKeywordsServer() WebServer {
	return &CognitubeAIServer{
		keywordsService: keywords.NewCognitubeKeywordsService(),
		concatService:   concat.NewConcatService(),
	}
}
