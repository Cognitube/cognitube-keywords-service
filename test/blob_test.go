package test

import (
	"cognitube.com/keywords-service/azure"
	"fmt"
	"testing"
)

func TestBlob_Update(t *testing.T) {
	client := azure.NewBlobClient()

	url, _ := client.UploadBlob("test-container", "test.txt", []byte("Hello, World!"))
	fmt.Println(url)
}
