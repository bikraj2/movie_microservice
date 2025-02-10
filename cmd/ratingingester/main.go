package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"bikraj.movie_microservice.net/rating/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	fmt.Println("Creating  A Kafka Producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.server": "localhost"})

	if err != nil {
		panic(err)
	}

	defer producer.Close()

	const fileName = "ratingsdata.json"

	fmt.Println("Reading rating events from file " + fileName)

	ratingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "raitings"

	if err := produceRatingEvents(topic, producer, ratingEvents); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second

	fmt.Printf("Wating %s until all event gets published", timeout.String())
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}

func produceRatingEvents(topic string, producer *kafka.Producer, events []model.RatingEvent) error {
	for _, ratingEvent := range events {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}

		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(encodedEvent),
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
