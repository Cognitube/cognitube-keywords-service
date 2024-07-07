package server

import (
	"fmt"
	"net/http"
)

type WebServer interface {
	StartListening(port string)
}

type CongitubeKeywordsServer struct {
	keywordsService ICognitubeKeywordsService
}

func (c *CongitubeKeywordsServer) StartListening(port string) {
	httpHandler := (http.HandlerFunc)(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		file, handler, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Error retrieving the audio file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Receive File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		keywords, err := c.keywordsService.GetKeyDescFromHttpAudioFile(file)
		if err != nil {
			http.Error(w, "Error getting keywords: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Done processing the audio file ", handler.Filename, " to keywords")
		w.Write([]byte(keywords))
	})

	http.HandleFunc("/api/getKeywords", httpHandler)

	fmt.Println("Server started at port ", port)
	http.ListenAndServe(":"+port, nil)
}

func NewCongitubeKeywordsServer() WebServer {
	return &CongitubeKeywordsServer{
		keywordsService: NewCognitubeKeywordsService(),
	}
}
