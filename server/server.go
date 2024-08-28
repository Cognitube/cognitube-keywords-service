package server

import (
	"cognitube.com/keywords-service/server/service/concat"
	"cognitube.com/keywords-service/server/service/keywords"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type WebServer interface {
	StartListening(port string)
}

type CongitubeAIServer struct {
	keywordsService keywords.ICognitubeKeywordsService
	concatService   concat.IConcatService
}

func (c *CongitubeAIServer) StartListening(port string) {
	router := mux.NewRouter()

	SetupRoutes(router, c)

	log.Println("Server started at port ", port)

	http.ListenAndServe(":"+port, router)
}

func NewCongitubeKeywordsServer() WebServer {
	return &CongitubeAIServer{
		keywordsService: keywords.NewCognitubeKeywordsService(),
		concatService:   concat.NewConcatService(),
	}
}
