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


## Optional

Use a kafka visualizer tool. I use [this one](https://github.com/manasb-uoe/kafka-visualizer/) 

```
java -jar ~/wednesday/kafka-visualizer/rest/target/rest-1.0-SNAPSHOT.jar --zookeeper=localhost:2181 --kafka=localhost:9092 --env=UAT
```