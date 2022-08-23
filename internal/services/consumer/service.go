package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/weezqyd/ringier/internal/adapters/msb"
	"github.com/weezqyd/ringier/internal/common/messages"
	"github.com/weezqyd/ringier/internal/services/dispatcher/grpc/dispatcher"
	"github.com/weezqyd/ringier/internal/services/persist/grpc/persister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

type Service struct {
	subscriber *msb.Subscriber
	config     *ServiceConfig
	persister  *grpc.ClientConn
	dispatcher *grpc.ClientConn
}
type ServiceConfig struct {
	RabbitMQURI      string
	ListenPort       uint
	SubscriptionName string
}

type ServiceOption func(service *Service) error

// NewService creates a new consumer Service instance.
func NewService(opts ...ServiceOption) (*Service, error) {
	service := &Service{}
	for _, fn := range opts {
		if err := fn(service); err != nil {
			return nil, err
		}
	}

	return service, nil
}

func WithConfig(config *ServiceConfig) ServiceOption {
	return func(service *Service) error {
		service.config = config
		return nil
	}
}

func WithSubscriber() ServiceOption {
	return func(service *Service) error {
		if service.config == nil {
			return fmt.Errorf("service config is missing ... ")
		}
		subscriber, err := msb.NewSubscriber(
			msb.WithRabbitSubscriber(service.config.RabbitMQURI),
		)
		if err != nil {
			return err
		}
		service.subscriber = subscriber
		return nil
	}
}

func WithPersister() ServiceOption {
	return func(service *Service) error {
		conn, err := grpc.Dial("persister:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("could not connect to persister service, %s", err)
		}
		service.persister = conn
		return nil
	}
}

func WithDispatcher() ServiceOption {
	return func(service *Service) error {
		conn, err := grpc.Dial("dispatcher:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("[dispatcher] could not connect to persister service, %s", err)
		}
		service.dispatcher = conn
		return nil
	}
}

func (s *Service) Start(ctx context.Context) error {
	return s.subscriber.Subscribe(ctx, s.config.SubscriptionName, func(msg *message.Message) error {
		req := &messages.EventRequest{}

		if err := json.Unmarshal(msg.Payload, req); err != nil {
			return fmt.Errorf("[consumer] could not decode event object %s", err)
		}
		if err := callService(ctx, s, "dispatcher", req); err != nil {
			return err
		}
		if err := callService(ctx, s, "persister", req); err != nil {
			return err
		}
		return nil
	})
}

func callService(ctx context.Context, s *Service, service string, req *messages.EventRequest) error {
	var err error
	res := &messages.EventResponse{}

	switch service {
	case "dispatcher":

		res, err = dispatcher.NewDispatcherServiceClient(s.dispatcher).Dispatch(ctx, req)
	case "persister":
		res, err = persister.NewPersistServiceClient(s.persister).Save(ctx, req)
	}
	if err != nil {
		return err
	}
	if res.Success != true {
		return fmt.Errorf("%s", strings.Join(res.Errors, ", "))
	}
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	if err := s.subscriber.Close(); err != nil {
		log.Printf("[consumer] could not close subscriber, %s", err)
	}
	if err := s.persister.Close(); err != nil {
		log.Printf("[consumer] could not close persister gRPC connection, %s", err)
	}
	if err := s.dispatcher.Close(); err != nil {
		log.Printf("[consumer] could not close dispatcher gRPC connection, %s", err)
	}

	return nil
}
