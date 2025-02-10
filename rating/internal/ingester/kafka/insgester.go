package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"bikraj.movie_microservice.net/rating/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.server":  addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Ingester{consumer, topic}, nil
}

func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("Consumer erro: %s", err.Error())
				continue
			}
			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("UnMarshal error: %s", err.Error())
				continue
			}
			ch <- event
		}
	}()
	return ch, nil
}
