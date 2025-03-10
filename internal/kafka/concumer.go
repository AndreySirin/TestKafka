package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func Consumer(address []string, ctx context.Context) error {
	kafkaConsumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		return fmt.Errorf("failed to start Kafka consumer: %v", err)
	}
	defer func() {
		errClose := kafkaConsumer.Close()
		if errClose != nil {
			log.Printf("failed to close Kafka consumer: %v", errClose)
		}
	}()

	pc, err := kafkaConsumer.ConsumePartition("messages", 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}
	defer func() {
		errClose := pc.Close()
		if errClose != nil {
			log.Printf("failed to close object PartitionConsumer: %v", errClose)
		}
	}()

	for {
		select {
		case _, ok := <-pc.Messages():
			if !ok {
				return fmt.Errorf("kafka consumer is closed")
			}

		case err = <-pc.Errors():
			log.Println("Kafka error:", err)

		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return nil
		}
	}
}
