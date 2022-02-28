package kafkaservice

import (
	"context"
	"fmt"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

type KAFKA_TOPIC string

const (
	ISSUED_TOKEN     KAFKA_TOPIC = "issued-token"
	NEW_USER_CREATED KAFKA_TOPIC = "new-user-created"
)

func getBrokers(count int) []string {
	broker := os.Getenv(fmt.Sprintf("KAFKA_HOST_%d", count))
	if broker == "" {
		return []string{}
	}
	return append([]string{broker}, getBrokers(count+1)...)
}

var tokenWriter *kafka.Writer = kafka.NewWriter(kafka.WriterConfig{
	Brokers:  getBrokers(1),
	Topic:    string(ISSUED_TOKEN),
	Balancer: &kafka.Hash{},
})

var newUserWriter *kafka.Writer = kafka.NewWriter(kafka.WriterConfig{
	Brokers:  getBrokers(1),
	Topic:    string(NEW_USER_CREATED),
	Balancer: &kafka.RoundRobin{},
})

func Produce(ctx context.Context, topic KAFKA_TOPIC, key []byte, value []byte) {
	var kafkaWriter *kafka.Writer
	switch topic {
	case ISSUED_TOKEN:
		kafkaWriter = tokenWriter
	case NEW_USER_CREATED:
		kafkaWriter = newUserWriter
	}

	if kafkaWriter == nil {
		fmt.Printf("Invalid topic: %s", topic)
		return
	}
	err := kafkaWriter.WriteMessages(ctx, kafka.Message{
		// create an arbitrary message payload for the value
		Value: value,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}

	// log a confirmation once the message is written
	fmt.Printf("::published \ntopic: %s\nkey:%s\nvalue:%s", topic, key, value)

}
