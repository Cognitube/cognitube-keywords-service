package server

import (
	"fmt"
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

	fmt.Println("Server started at port ", port)

	http.ListenAndServe(":"+port, router)
}

func NewCongitubeKeywordsServer() WebServer {
	return &CongitubeKeywordsServer{
		keywordsService: NewCognitubeKeywordsService(),
	}
}
