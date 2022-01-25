package kafkaservice

import (
	"context"
	"fmt"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

// the topic and broker address are initialized as constants
const (
	topic = "message-log"
)

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func Initiate() {
	ctx := context.Background()
	go consume(ctx)
}
