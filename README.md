# go-kafka-example

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
