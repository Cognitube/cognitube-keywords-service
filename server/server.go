package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type WebServer interface {
	StartListening(port string)
}

type CongitubeKeywordsServer struct {
	keywordsService ICognitubeKeywordsService
}

func (c *CongitubeKeywordsServer) StartListening(port string) {
	router := mux.NewRouter()

	SetupRoutes(router, c.keywordsService)

	log.Println("Server started at port ", port)

	http.ListenAndServe(":"+port, router)
}

func NewCongitubeKeywordsServer() WebServer {
	return &CongitubeKeywordsServer{
		keywordsService: NewCognitubeKeywordsService(),
	}
}
