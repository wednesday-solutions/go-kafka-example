# go-kafka-example


## Prerequisites


- docker
- zookeeper
- kafka
- aws copilot 

This is a monorepo setup with support for CI/CD. The applications in the monorepo are based on [go-template](https://github.com/wednesday-solutions/go-template)

The [producer](./producer) publishes messages to 2 kafka topics - issued-token, new-user-created. The [consumer](./consumer) consumes these messages. This is a working setup for message broking in golang using kafka in a micro-service environment.

The consumer also exposes an API that makes an inter-service API call to the producer to fulfil the request. 

This application is deployed on AWS ECS using AWS Copilot. They are deployed as 2 services in a cluster - hence are able to use the the service discover DNS for interservice communication. 

## Run the application

1. Start zookeeper
        ```
        zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties
        ```
2. Start the kafka service
        ```
        kafka-server-start /usr/local/etc/kafka/server.properties
        ```
3. Start the consumer 
        ```
        cd consumer 
        go run cmd/server/main.go
        ```
4. Start the producer
        ```
        cd ../producer
        go run cmd/server/main.go
        ```


## Optional

Use a kafka visualizer tool. I use [this one](https://github.com/manasb-uoe/kafka-visualizer/) 

```
java -jar ~/wednesday/kafka-visualizer/rest/target/rest-1.0-SNAPSHOT.jar --zookeeper=localhost:2181 --kafka=localhost:9092 --env=UAT
```

## Kafka cli commands

- descirbe all consumers, topics and partition detiails
```
kafka-consumer-groups --bootstrap-server localhost:9092 --all-groups --describe
```

## Config

The power of kafka is seen at scale which is possible because we can have multiple partitions. 
Use the config files found in [docs/config](./docs/config/) to leverage it to some extent. I have 8 partitions set up.

Here is a small write up to explain how kafka works

- in kafka the main queue is broken down into many subqueues. Each of these subqueues is called a partition.
- A server that holds one or more partition is called a broker. Each item in the partition is called a record.
- The field that decides which partition the record will be stored in is called a key. [#](./producer/pkg/utl/kafkaservice/kafka.go#L21)
- If no key is specified then a random partition is assigned.[#](./producer/pkg/utl/kafkaservice/kafka.go#L27)
- A group of partitions handling the same kind of data is called a topic. [#](./producer/pkg/utl/kafkaservice/kafka.go#L43)
- An offset is a sequential number provided to each record.
- A record in a topic is identified by a partition number and an offset. [#](./consumer/pkg/utl/kafkaservice/kafka.go#L51)
- Having one consumer per partition guarantees ordering per game. Consumers can be scaled easily and without a lot of performance or cost impact.
- This is because kafka only needs to maintain the latest offset read by each consumer. [here we're letting kafka handle committing](./consumer/pkg/utl/kafkaservice/kafka.go#L70)
- Typically consumers read one record at a time, and pickup where they left off after a restart. [here we're manually handling committing](./consumer/pkg/utl/kafkaservice/kafka.go#L35)
- It's  quite common to have consumers read all the records from the beginning on startup
- Consumers in a consumer group do not share partitions. Each consumer would read different records from the other consumers.
- Multiple consumer groups are useful when you have different applications reading the same content
- Kafka has retention policies. For example after 24 hours the kafka queue will be cleaned.
- Kafka can also store all records on a persistent storage. This makes it fault-tolerant and durable. So if the broker goes down it can recover when it comes back up.
- Replication Factor: Kafka replicates partitions so when a broker goes down, a backup parition takes over and processing can resume. This is configured using the replication factor. if you have 3, it means you have 3 copies of a partition. 1 leader and 2 backups. This means we can tolerate upto 2 brokers going down at the same time.