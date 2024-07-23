package test

import (
	"cognitube.com/keywords-service/publish"
	"cognitube.com/keywords-service/result"
	"context"
	"github.com/segmentio/kafka-go"
	"os"
	"testing"
)

func TestKafkaPublisher_Publish(t *testing.T) {
	k := publish.NewKafkaPublisher("localhost", "9092")
	k.Publish("test", []byte("test message"))
	// get the message from kafka and assert it
	// using kafka-go
}

func TestTranscriptionResultPublisher_PublishTranscriptionResult(t *testing.T) {
	os.Setenv("APPLICATION_KAFKA_HOST", "localhost")
	os.Setenv("APPLICATION_KAFKA_PORT", "9092")
	os.Setenv("APPLICATION_KAFKA_TOPIC", "test")

	p := publish.NewKafkaTranscriptionPublisher()
	err := p.PublishTranscriptionResult(&result.Result{
		Success:       true,
		VideoID:       "test",
		TranscriptURL: "test url to transcript",
		KeywordsURL:   "test url to keywords",
	})
	if err != nil {
		t.Error(err)
		return
	}
	// get the message from kafka and assert it
	// using kafka-go
}

func TestConsumer(t *testing.T) {
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
