package mq

import (
	"encoding/json"

	"cognitube.com/keywords-service/env"
	"cognitube.com/keywords-service/server/service/keywords/result"
)

type TranscriptionPublisher interface {
	PublishTranscriptionResult(result *result.Result) error
	Publish(topic string, message []byte) error
}
type KafkaTranscriptionPublisher struct {
	KafkaPublisher
}

func (p *KafkaTranscriptionPublisher) PublishTranscriptionResult(result *result.Result) error {
	msg, _ := json.Marshal(result)
	topic := env.GetInstance().KafkaTopic
	return p.Publish(topic, msg)
}

func (p *KafkaTranscriptionPublisher) Publish(topic string, message []byte) error {
	return p.KafkaPublisher.Publish(topic, message)
}

func NewKafkaTranscriptionPublisher() TranscriptionPublisher {
	host := env.GetInstance().KafkaHost
	port := env.GetInstance().KafkaPort
	return &KafkaTranscriptionPublisher{
		KafkaPublisher: KafkaPublisher{
			Url: host + ":" + port,
		},
	}
}
