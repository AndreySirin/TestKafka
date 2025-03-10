package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	KafkaMessagesSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_sent_total",
			Help: "Total number of successfully sent Kafka messages.",
		},
	)

	KafkaMessagesErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_errors_total",
			Help: "Total number of Kafka messages that failed to send.",
		},
	)
)

func Init() {
	prometheus.MustRegister(KafkaMessagesSent)
	prometheus.MustRegister(KafkaMessagesErrors)
}
