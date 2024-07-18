package publish

import (
	"cognitube.com/keywords-service/env"
	"cognitube.com/keywords-service/transcription"
	"encoding/json"
)

type TranscriptionPublisher interface {
	PublishTranscriptionResult(result *transcription.TranscriptionResult) error
}
type KafkaTranscriptionPublisher struct {
	KafkaPublisher
}

func (p *KafkaTranscriptionPublisher) PublishTranscriptionResult(result *transcription.TranscriptionResult) error {
	msg, _ := json.Marshal(result)
	topic := env.GetInstance().KafkaTopic
	return p.Publish(topic, msg)
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
