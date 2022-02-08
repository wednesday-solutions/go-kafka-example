package kafkaservice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	kafka "github.com/segmentio/kafka-go"
	models "github.com/wednesday-solutions/go-template-consumer/models"
)

type KAFKA_TOPIC string

const (
	ISSUED_TOKEN     = "issued-token"
	NEW_USER_CREATED = "new-user-created"
)

func consumeIssuedToken(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	issuedToken := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   string(ISSUED_TOKEN),
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := issuedToken.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		var user models.User
		e := json.Unmarshal(msg.Value, &user)
		if e != nil {
			fmt.Print("error while unmarshalling", e)
		} else {
			fmt.Print("\n\nNew token issued by the following user:", user)

		}
	}
}

func consumeNewUserCreated(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	newUserCreatedReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   string(NEW_USER_CREATED),
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := newUserCreatedReader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		var user models.User
		e := json.Unmarshal(msg.Value, &user)
		if e != nil {
			fmt.Print("error while unmarshalling", e)
		} else {
			fmt.Print("\n\nNew user created:", user)
		}
	}
}

func Initiate() {
	ctx := context.Background()
	go consumeIssuedToken(ctx)
	go consumeNewUserCreated(ctx)
}
