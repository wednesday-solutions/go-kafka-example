package kafkaservice

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"

	rediscache "consumer/pkg/utl/rediscache"

	"github.com/gocraft/work"
	hashstructure "github.com/mitchellh/hashstructure/v2"
)

// Cont ...
type Cont struct {
}

var GO_KAFKA_EXAMPLE = "go-kafka-example"
var PRUNE_CACHE = "prune_cache"

var SCHEDULED_KAFKA_RETRY = "scheduled_kafka_retry"
var SCHEDULED_CACHE_PRUNE = "scheduled_cache_prune"

var ExponentialBackOffExp int64 = 5
var ExponentialBackOffInterval int64 = 60
var ExponentialBackOffMaxRetry = 5

func getExponentialBackOff(count int) int64 {
	count = count - 1
	if count < 0 {
		return 0
	} else if count == 0 {
		return ExponentialBackOffInterval
	}
	exp := math.Pow(float64(ExponentialBackOffExp), float64(count))
	return ExponentialBackOffInterval*int64(exp) + getExponentialBackOff(count)
}

func InitScheduledJobs() {
	pool := work.NewWorkerPool(Cont{}, 10, GO_KAFKA_EXAMPLE, rediscache.RedisPool)
	fmt.Println("-----------InitScheduledJobs---------------")
	pool.Job(SCHEDULED_KAFKA_RETRY, KakfaRetryProcessor)
	pool.Job(SCHEDULED_CACHE_PRUNE, CachePruningProcessor)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func KakfaRetryProcessor(job *work.Job) error {
	byteRretryMessage, err := json.Marshal(job.Args["retryMessage"])
	if err != nil {
		fmt.Println("there was an error while marshaling bro")
	}
	var retryMessage RetryMessage
	err = json.Unmarshal(byteRretryMessage, &retryMessage)
	if err != nil {
		fmt.Println("there was an error while unmarshaling bro")
	}
	fmt.Println("\nrepublishing message to kafka ", retryMessage)
	Produce(
		context.Background(),
		KAFKA_TOPIC(retryMessage.Topic),
		[]byte(retryMessage.Key),
		retryMessage.Message,
	)
	return err
}

func CachePruningProcessor(job *work.Job) error {
	hash, _ := job.Args["hash"].(string)
	a, err := rediscache.RemoveKeyValue(hash)
	fmt.Println("a,err", a, err)
	return nil
}

func ScheduleCachePruning(hash string) {
	enqueuer := work.NewEnqueuer(GO_KAFKA_EXAMPLE, rediscache.RedisPool)
	secondsFromNow := getExponentialBackOff(ExponentialBackOffMaxRetry + 1)
	fmt.Println("ScheduleCachePruning for hash:", hash, " ", secondsFromNow, "seconds from now")
	job, err := enqueuer.EnqueueIn(SCHEDULED_CACHE_PRUNE, secondsFromNow, map[string]interface{}{"hash": hash})
	if err != nil {
		fmt.Println("Error in enqueuing job for cache pruning ", err)
	} else {
		var name string = "couldn't get name"
		if job != nil {
			name = job.Name
		}
		fmt.Println("\n---------------Job enqueued: ", name, "--------------------")
	}
}

func ScheduleKafkaRetry(retryMessage RetryMessage) {
	hash, err := hashstructure.Hash(map[string]string{
		"topic":   retryMessage.Topic,
		"message": string(retryMessage.Message),
	}, hashstructure.FormatV2, nil)
	if err != nil {
		// send slack alert
		fmt.Println("error while creating a hash in the ScheduleKafkaRetry function.\n", err)
	}
	count := 1
	val, err := rediscache.GetKeyValue(fmt.Sprint(hash))
	if err != nil {
		// send slack alert
		fmt.Println("error while getting value from redis in the ScheduleKafkaRetry function.\n", err)
	}
	if val == nil {
		ScheduleCachePruning(fmt.Sprint(hash))
	} else {
		b := val.([]byte)
		err = json.Unmarshal(b, &count)
		if err != nil {
			// send slack alert
			fmt.Println("error while unmarshaling count from redis in the ScheduleKafkaRetry function.\n", err)
		}
		count = count + 1
	}
	if count < ExponentialBackOffMaxRetry {
		err := rediscache.SetKeyValue(fmt.Sprint(hash), count)
		if err != nil {
			// send slack alert
			fmt.Println("error while udpating the count in redis", err)
		}
	} else {
		// sending items to the dead letter queue
		go Produce(context.Background(), DEAD_LETTER_QUEUE, []byte(retryMessage.Key),
			retryMessage.Message)
	}
	enqueuer := work.NewEnqueuer(GO_KAFKA_EXAMPLE, rediscache.RedisPool)
	secondsFromNow := getExponentialBackOff(count)
	fmt.Println("\nScheduling kafka retry in: ", secondsFromNow, " seconds from now. Topic:",
		retryMessage.Topic, "Count: ", count)
	job, err := enqueuer.EnqueueIn(SCHEDULED_KAFKA_RETRY, secondsFromNow,
		map[string]interface{}{"retryMessage": retryMessage})
	if err != nil {
		fmt.Println("Error in enqueuing job for kafka retry ", err)
	} else {
		var name string = "couldn't get name"
		if job != nil {
			name = job.Name
		}
		fmt.Println("\n---------------Job enqueued: ", name, "--------------------")
	}
}
