package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Kafka interface {
	Send(string, string) error
}

type KaFka struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*KaFka, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("producer creation error: %w", err)
	}
	return &KaFka{
		producer: producer,
	}, nil
}

func (p *KaFka) Send(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		KafkaMessagesErrors.Inc()
		return fmt.Errorf("error sending the message: %v", err)
	}
	KafkaMessagesSent.Inc()
	return nil
}
