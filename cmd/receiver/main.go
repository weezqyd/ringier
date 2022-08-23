package main

import (
	"github.com/gin-gonic/gin"
	"github.com/weezqyd/ringier/internal/services/receiver"
	"log"
)

func main() {
	config := &receiver.ServiceConfig{
		ListenPort:   9000,
		RabbitMQURI:  "amqp://guest:guest@rabbitmq:5672/",
		KafkaBrokers: []string{"kafka:9092"},
		Topic:        "ringier.events.users",
	}
	service, err := receiver.NewReceiverService(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	service.Start(gin.Default())
}
