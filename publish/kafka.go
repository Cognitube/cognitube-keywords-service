package publish

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Publisher interface {
	Publish(topic string, message []byte) error
}

type KafkaPublisher struct {
	Url string
}

func (p *KafkaPublisher) Publish(topic string, message []byte) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", p.Url, topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Write(
		message,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	return nil
}

func NewKafkaPublisher(host, port string) Publisher {
	return &KafkaPublisher{host + ":" + port}
}
