package kafkaservice

import (
	"context"
	"encoding/json"
	"fmt"

	"os"

	models "consumer/models"
	"consumer/pkg/utl/convert"
	"consumer/resolver"

	kafka "github.com/segmentio/kafka-go"
)

type KAFKA_TOPIC string

const (
	ISSUED_TOKEN     = "issued-token"
	NEW_USER_CREATED = "new-user-created"
)

func getBrokers(count int) []string {
	broker := os.Getenv(fmt.Sprintf("KAFKA_HOST_%d", count))
	if broker == "" {
		return []string{}
	}
	return append([]string{broker}, getBrokers(count+1)...)
}

func consumeIssuedToken(ctx context.Context, r *resolver.Resolver) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	issuedToken := kafka.NewReader(kafka.ReaderConfig{
		Brokers: getBrokers(1),
		Topic:   string(ISSUED_TOKEN),
		GroupID: "my-group",
	})
	for {
		// here we will manual commit messages using FetchMessage instead of read message
		msg, err := issuedToken.FetchMessage(ctx)
		if err != nil {
			fmt.Print("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		var user models.User
		e := json.Unmarshal(msg.Value, &user)
		if e != nil {
			fmt.Print("error while unmarshalling", e)
		} else {
			fmt.Printf("\n\nNew token issued by user:%d, partition: %d, offset: %d", user.ID, msg.Partition, msg.Offset)
			for _, o := range r.Observers {
				o <- convert.UserToGraphQlUser(&user)
			}
			if err := issuedToken.CommitMessages(ctx, msg); err != nil {
				fmt.Printf("failed to commit message. "+
					"Partition: %d, offset: %d", msg.Partition, msg.Offset)
			} else {
				fmt.Print("\n\ncommited message")
			}
		}
	}
}

func consumeNewUserCreated(ctx context.Context, r *resolver.Resolver) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	newUserCreatedReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: getBrokers(1),
		Topic:   string(NEW_USER_CREATED),
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := newUserCreatedReader.ReadMessage(ctx)
		if err != nil {
			fmt.Print("\n\ncould not read message " + err.Error())
		}
		var user models.User
		e := json.Unmarshal(msg.Value, &user)
		if e != nil {
			fmt.Print("\n\nerror while unmarshalling", e)
		} else {
			fmt.Printf("\n\nNew user created: %d, partition: %d, offset: %d", user.ID, msg.Partition, msg.Offset)
			for _, o := range r.Observers {
				o <- convert.UserToGraphQlUser(&user)
			}
		}
	}
}

func Initiate(r *resolver.Resolver) {
	ctx := context.Background()
	go consumeIssuedToken(ctx, r)
	go consumeNewUserCreated(ctx, r)
}
