# Consumer

This application consumes messages from 2 topcis - [issued-token](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/pkg/utl/kafkaservice/kafka.go#L36) and [new-user-created](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/pkg/utl/kafkaservice/kafka.go#L71) 

The reader responsible for reading messages on the "issued-token" topic requires manual committing of messages and hence used the FetchMessage API coupled with the CommitMessages API from kafka-go. [#](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/pkg/utl/kafkaservice/kafka.go#L43)

However in case of the "new-user-created" topic committing of messages is delegated to kafka-go itself using the ReadMessage API. [#](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/pkg/utl/kafkaservice/kafka.go#L78)

The application comes with out of the box support for graphql subscriptions. For the purpose of demoing the flow of orchestrating subscritpions on consumption of messages the userNotification subscritpion is triggered when the a new token is issued or a new user is created. 

- [new-user-created](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/resolver/user_mutations.resolvers.go#L47-L51)
- [issued-token](https://github.com/wednesday-solutions/go-kafka-example/blob/main/consumer/pkg/utl/kafkaservice/kafka.go#L54-L56)

You can use the following subscription to receive notifications. 

```

  subscription UserNotification {
    userNotification {
      id
      firstName
      lastName
    }
  }
```
