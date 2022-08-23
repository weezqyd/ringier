package msb

import (
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
)

// PublisherOption is a function for setting a Publisher value.
type PublisherOption func(service *Publisher) error

type Publisher struct {
	kafkaPublisher  *kafka.Publisher
	rabbitPublisher *amqp.Publisher
	topic           string
}

// NewPublisher creates a new Publisher instance with multiple brokers.
func NewPublisher(opts ...PublisherOption) (*Publisher, error) {
	publisher := &Publisher{}
	for _, fn := range opts {
		if err := fn(publisher); err != nil {
			return nil, err
		}
	}

	return publisher, nil
}

func WithRabbitPublisher(uri string) PublisherOption {
	return func(publisher *Publisher) error {
		var err error
		if publisher.rabbitPublisher, err = NewRabbitPublisher(uri); err != nil {
			return err
		}
		return nil
	}
}

func WithTopic(topic string) PublisherOption {
	return func(publisher *Publisher) error {
		if publisher.rabbitPublisher != nil {
			channel, err := publisher.rabbitPublisher.Connection().Channel()
			if err != nil {
				return fmt.Errorf("could not get channel from connection err: %s", err.Error())
			}
			defer channel.Close()
			_, err = channel.QueueDeclare(topic, true, false, false, false, nil)
			if err != nil {
				return fmt.Errorf("could not declare queue err: %s", err.Error())
			}
		}
		publisher.topic = topic
		return nil
	}
}

func WithKafkaPublisher(brokers []string) PublisherOption {
	return func(publisher *Publisher) error {
		var err error
		if publisher.kafkaPublisher, err = NewKafkaPublisher(brokers); err != nil {
			return err
		}
		return nil
	}
}

func (p *Publisher) Publish(message *message.Message) error {
	var err error
	if p.rabbitPublisher != nil && p.rabbitPublisher.IsConnected() {
		// try to publish to rabbit the primary MSB first
		err = p.rabbitPublisher.Publish(p.topic, message)
		if err == nil {
			return nil
		}
		log.Printf("[WARN] Failed to publish to rabbit: %s", err)
	}
	if p.kafkaPublisher != nil {
		err = p.kafkaPublisher.Publish(p.topic, message)
		if err == nil {
			return nil
		}
		log.Printf("[WARN] Failed to publish to kafka: %s", err)
	}
	return errors.New("no publisher available at the moment, try again later")
}
