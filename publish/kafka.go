package publish

import (
	"cognitube.com/keywords-service/env"
	"context"
	"crypto/tls"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"strings"
)

type Publisher interface {
	Publish(topic string, message []byte) error
}

type KafkaPublisher struct {
	Url string
}

func (p *KafkaPublisher) PublishProd(topic string, message []byte) error {
	eventHubNamespace := env.GetInstance().EventHubNamespace
	eventHubName := env.GetInstance().EventHubName
	connectionString := env.GetInstance().EventHubConnectionString
	username, password := parseConnectionString(connectionString)

	// Set up SASL configuration
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(eventHubNamespace),
		Topic:    eventHubName,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello, Azure Event Hubs!"),
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (p *KafkaPublisher) Publish(topic string, message []byte) error {
	log.Println("Try publish message to Kafka: " + string(message))
	if !env.GetInstance().Debug {
		return p.PublishProd(topic, message)
	}

	w := &kafka.Writer{
		Addr:  kafka.TCP(p.Url),
		Topic: topic,
	}
	defer w.Close()
	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})

	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func NewKafkaPublisher(host, port string) Publisher {
	return &KafkaPublisher{host + ":" + port}
}

func parseConnectionString(connectionString string) (string, string) {
	// The connection string format is: Endpoint=sb://<NAMESPACE>.servicebus.windows.net/;SharedAccessKeyName=<KEY_NAME>;SharedAccessKey=<KEY_VALUE>
	// Extract SharedAccessKeyName and SharedAccessKey from the connection string
	var username, password string
	for _, part := range strings.Split(connectionString, ";") {
		if strings.HasPrefix(part, "SharedAccessKeyName=") {
			username = strings.TrimPrefix(part, "SharedAccessKeyName=")
		} else if strings.HasPrefix(part, "SharedAccessKey=") {
			password = strings.TrimPrefix(part, "SharedAccessKey=")
		}
	}
	return username, password
}
