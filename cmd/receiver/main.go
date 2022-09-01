package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weezqyd/ringier/internal/services/receiver"
	"log"
)

var totalReqs = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "receiver_http_requests_total",
		Help: "How many requests are processed, partitioned by event and status code",
	},
	[]string{"event", "code"},
)

var latency = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "receiver_latency_time_seconds",
		Help: "Duration taken to handle the http request.",
	},
	[]string{"event"},
)

func main() {
	config := &receiver.ServiceConfig{
		ListenPort:   9000,
		RabbitMQURI:  "amqp://guest:guest@rabbitmq:5672/",
		KafkaBrokers: []string{"kafka:9092"},
		Topic:        "ringier.events.users",
	}
	service, err := receiver.NewService(
		receiver.WithConfig(config),
		receiver.WithPublisher(),
		receiver.WithMetrics(&receiver.Metrics{
			Latency:       latency,
			TotalRequests: totalReqs,
		}),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	app := gin.Default()

	handler := promhttp.Handler()
	app.GET("/metrics", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})

	service.Start(app)
}
