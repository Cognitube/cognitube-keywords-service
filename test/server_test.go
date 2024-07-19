package test

import (
	"context"
	"github.com/carlmjohnson/requests"
	"github.com/segmentio/kafka-go"
	"net/http"
	"testing"
)

// Step 0: Start the server by running the main.go file with environment variables set

func RegisterCallback(t *testing.T) {
	// This will send a POST request to the server to register a callback
	var body = struct {
		Url string `json:"url"`
	}{
		Url: "http://localhost:8848/static/test.mp3",
	}
	err := requests.URL("http://localhost:8082/api/v1/callback").
		Method("POST").
		BodyJSON(body).
		Fetch(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
}

// Step 1: Start a file server on port 8848 using the static directory to /static
func TestFileServer(t *testing.T) {
	// This will start a file server on port 8848 using the static directory to /static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8848", nil)
}

// Step 2: Listen kafka messages to display the final result
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

// Step 3: send a POST request to the server to create a transcription
func TestAll(t *testing.T) {
	var body = struct {
		AudioURL    string `json:"audio_url"`
		DisplayName string `json:"display_name"`
	}{
		AudioURL:    "http://localhost:8848/static/test.mp3",
		DisplayName: "test",
	}
	var id string
	err := requests.URL("http://localhost:8082/api/v1/transcription/create").
		Method("POST").
		BodyJSON(body).
		ToString(&id).
		Fetch(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)
}
