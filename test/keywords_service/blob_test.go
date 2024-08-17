package keywords_service

import (
	"cognitube.com/keywords-service/dal/azure"
	"log"
	"testing"
)

func TestBlob_Update(t *testing.T) {
	client := azure.NewBlobClient()

	url, _ := client.UploadBlob("test-container", "test.txt", []byte("Hello, World!"))
	log.Println(url)
}
