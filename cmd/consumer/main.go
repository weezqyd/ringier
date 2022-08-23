package main

import (
	"context"
	"github.com/weezqyd/ringier/internal/services/consumer"
	"log"
)

func main() {
	config := &consumer.ServiceConfig{
		RabbitMQURI:      "amqp://guest:guest@rabbitmq:5672/",
		SubscriptionName: "ringier.events.users",
		ListenPort:       9001,
	}
	ctx := context.Background()
	service, err := consumer.NewService(
		consumer.WithConfig(config),
		consumer.WithSubscriber(),
		consumer.WithPersister(),
		consumer.WithDispatcher(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop(ctx)
	log.Fatal(service.Start(ctx))
}
