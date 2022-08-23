package msb

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

func NewRabbitPublisher(uri string) (*amqp.Publisher, error) {
	return amqp.NewPublisher(
		amqp.NewDurableQueueConfig(uri),
		watermill.NewStdLogger(false, false),
	)
}

func NewKafkaPublisher(brokers []string) (*kafka.Publisher, error) {
	return kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
}

func NewRabbitSubscriber(uri string) (*amqp.Subscriber, error) {
	return amqp.NewSubscriber(
		amqp.NewDurableQueueConfig(uri),
		watermill.NewStdLogger(false, false),
	)
}
