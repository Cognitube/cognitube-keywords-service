package test

import (
	"cognitube.com/keywords-service/azure"
	"log"
	"testing"
)

func TestBlob_Update(t *testing.T) {
	client := azure.NewBlobClient()

	url, _ := client.UploadBlob("test-container", "test.txt", []byte("Hello, World!"))
	log.Println(url)
}
