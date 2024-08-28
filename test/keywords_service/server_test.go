package keywords_service

import (
	"context"
	"github.com/carlmjohnson/requests"
	"github.com/segmentio/kafka-go"
	"testing"
)

// Step 0: Start the server by compile and run main.go
// Step 1: Listen kafka messages to display the final result
func TestKafkaConsumer(t *testing.T) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "test-group",
		Topic:     "test",
		MaxBytes:  10e6, // 10KB
		Partition: 0,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		t.Log(string(m.Value))
	}
	r.Close()
}

// Step 2: send a POST request to the server to create a transcription
func TestNoConversation(t *testing.T) {
	audioUrl := "https://cognitube.blob.core.windows.net/audio-container/noconversation.mp3"
	videoID := "test-non-conversation-video-id"
	id, err := RequestLocal(audioUrl, videoID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)
}

func TestShort(t *testing.T) {
	audioUrl := "https://cognitube.blob.core.windows.net/audio-container/1720453568685-audio_9cede306-a9ea-4733-b27c-7097f3e89a31.oga"
	videoID := "test-short-video-id"
	id, err := RequestLocal(audioUrl, videoID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)
}

func TestLong(t *testing.T) {
	audioUrl := "https://cognitube.blob.core.windows.net/audio-container/1720458538777-audio_419fa396-d881-466e-b327-df94de0fc926.oga"
	videoID := "test-long-video-id"
	id, err := RequestLocal(audioUrl, videoID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)

}

func RequestLocal(audioUrl, videoId string) (string, error) {
	var id string
	err := requests.URL("http://localhost:8082/api/v1/transcription/create").
		Method("POST").
		BodyJSON(struct {
			AudioURL string `json:"audioUrl"`
			VideoID  string `json:"videoId"`
		}{AudioURL: audioUrl, VideoID: videoId}).
		ToString(&id).
		Fetch(context.Background())
	return id, err
}

// Optional: simulate the callback from the server
func TestCallback(t *testing.T) {
	err := requests.URL("http://localhost:8082/api/v1/callback").
		Method("POST").
		BodyJSON(struct {
			Self string `json:"self"`
		}{Self: "https://eastus.api.cognitive.microsoft.com/speechtotext/v3.1/transcriptions/2fe33d35-edac-4f84-99dd-dfc6badadc09"}).
		Fetch(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
}
