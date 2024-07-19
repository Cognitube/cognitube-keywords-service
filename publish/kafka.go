package publish

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(topic string, message []byte) error
}

type KafkaPublisher struct {
	Url string
}

func (p *KafkaPublisher) Publish(topic string, message []byte) error {
	w := &kafka.Writer{
		Addr:  kafka.TCP(p.Url),
		Topic: topic,
	}
	defer w.Close()
	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
	return err
}

func NewKafkaPublisher(host, port string) Publisher {
	return &KafkaPublisher{host + ":" + port}
}
