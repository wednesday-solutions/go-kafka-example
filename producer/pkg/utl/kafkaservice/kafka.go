package kafkaservice

import (
	"context"
	"fmt"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

type KAFKA_TOPIC string

const (
	ISSUED_TOKEN         KAFKA_TOPIC = "issued-token"
	NEW_USER_CREATED     KAFKA_TOPIC = "new-user-created"
	SIDE_TOPIC_FOR_RETRY KAFKA_TOPIC = "side-topic-for-retry"
)

func getBrokers(count int) []string {
	broker := os.Getenv(fmt.Sprintf("KAFKA_HOST_%d", count))
	fmt.Print("KAFKA_HOST: ", fmt.Sprintf("KAFKA_HOST_%d", count), "broker:", broker)
	if broker == "" {
		return []string{}
	}
	return append([]string{broker}, getBrokers(count+1)...)
}

var tokenWriter *kafka.Writer
var newUserWriter *kafka.Writer
var sideTopicWriter *kafka.Writer

func Produce(ctx context.Context, topic KAFKA_TOPIC, key []byte, value []byte) {
	var kafkaWriter *kafka.Writer

	switch topic {
	case ISSUED_TOKEN:
		if tokenWriter == nil {
			tokenWriter = kafka.NewWriter(kafka.WriterConfig{
				Brokers:  getBrokers(1),
				Topic:    string(ISSUED_TOKEN),
				Balancer: &kafka.Hash{},
			})
		}
		kafkaWriter = tokenWriter
	case NEW_USER_CREATED:
		if newUserWriter == nil {
			newUserWriter = kafka.NewWriter(kafka.WriterConfig{
				Brokers:  getBrokers(1),
				Topic:    string(NEW_USER_CREATED),
				Balancer: &kafka.RoundRobin{},
			})
		}
		kafkaWriter = newUserWriter

	case SIDE_TOPIC_FOR_RETRY:
		if sideTopicWriter == nil {
			sideTopicWriter = kafka.NewWriter(kafka.WriterConfig{
				Brokers:  getBrokers(1),
				Topic:    string(SIDE_TOPIC_FOR_RETRY),
				Balancer: &kafka.RoundRobin{},
			})
		}
		kafkaWriter = newUserWriter
	}

	if kafkaWriter == nil {
		fmt.Printf("Invalid topic: %s", topic)
		return
	}
	fmt.Print("list of brokers:", getBrokers(1))
	err := kafkaWriter.WriteMessages(ctx, kafka.Message{
		// create an arbitrary message payload for the value
		Value: value,
	})
	if err != nil {
		fmt.Print("could not write message " + err.Error())
	}

	// log a confirmation once the message is written
	fmt.Printf("\n\n::published \ntopic: %s\nkey:%s\nvalue:%s", topic, key, value)

}
