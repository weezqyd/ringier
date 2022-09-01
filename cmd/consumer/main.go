package main

import (
	"context"
	"github.com/weezqyd/ringier/internal/services/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := &consumer.ServiceConfig{
		RabbitMQURI:      "amqp://guest:guest@rabbitmq:5672/",
		SubscriptionName: "ringier.events.users",
		ListenPort:       9001,
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
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
