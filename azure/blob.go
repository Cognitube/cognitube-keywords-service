package azure

import (
	"cognitube.com/keywords-service/env"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"log"
)

type IBlobClient interface {
	UploadBlob(container string, blobName string, data []byte) (url string, err error)
}

type BlobClient struct {
	_client *azblob.Client
}

func (b *BlobClient) UploadBlob(container string, blobName string, data []byte) (string, error) {
	_, err := b._client.UploadBuffer(context.TODO(), container, blobName, data, nil)
	if err != nil {
		log.Println("Failed to upload blob: ", err.Error())
		return "", err
	}
	return fmt.Sprintf("%s%s/%s", b._client.URL(), container, blobName), err
}

func (b *BlobClient) UploadTranscript(filename, transcript string) (string, error) {
	return b.UploadBlob(env.GetInstance().TranscriptContainerName, filename, []byte(transcript))
}

func (b *BlobClient) UploadKeywords(filename string, desc string) (string, error) {
	return b.UploadBlob(env.GetInstance().KeywordsContainerName, filename, []byte(desc))
}

func NewBlobClient() *BlobClient {
	client, _ := azblob.NewClientFromConnectionString(env.GetInstance().BlobConnectString, nil)

	return &BlobClient{
		_client: client,
	}
}
