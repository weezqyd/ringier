package msb

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/net/context"
	"log"
)

type Subscriber struct {
	conn *amqp.Subscriber
}

type SubscriberFunc func(msg *message.Message) error

// SubscriberOption is a function for setting a Subscriber value.
type SubscriberOption func(service *Subscriber) error

// NewSubscriber creates a new Subscriber instance with multiple brokers.
func NewSubscriber(opts ...SubscriberOption) (*Subscriber, error) {
	subscriber := &Subscriber{}
	for _, fn := range opts {
		if err := fn(subscriber); err != nil {
			return nil, err
		}
	}

	return subscriber, nil
}

func WithRabbitSubscriber(uri string) SubscriberOption {
	return func(subscriber *Subscriber) error {
		var err error
		if subscriber.conn, err = NewRabbitSubscriber(uri); err != nil {
			return err
		}
		return nil
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, queue string, subscriberFunc SubscriberFunc) error {
	d, err := s.conn.Subscribe(ctx, queue)
	if err != nil {
		return fmt.Errorf("could not subscribe to topic[ %s] error: %s", queue, err)
	}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("subscriber exiting error: %s", ctx.Err())

		case msg := <-d:
			go func(msg *message.Message) {
				log.Printf("[consumer] new message => %s", msg.Payload)
				if err := subscriberFunc(msg); err != nil {
					log.Printf("[consumer] rejecting message because of error :%s", err)
					//if we encounter a message in our subscriber func we will Nack the message
					msg.Nack()
				}
			}(msg)
		}
	}
}

func (s *Subscriber) Close() error {
	return s.conn.Close()
}
