package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weezqyd/ringier/internal/adapters/database"
	gRPC "github.com/weezqyd/ringier/internal/services/persist/grpc"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

var persisterTotalReqs = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "persister_processed_total",
		Help: "How many events are processed, partitioned by event and status.",
	},
	[]string{"event", "status"},
)

var persisterLatency = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "persister_latency_time_seconds",
		Help: "Duration taken to handle save event.",
	},
	[]string{"event"},
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()
	db, err := database.NewMongoClient(
		ctx,
		"mongodb://root:secret@mongodb:27017/",
		"ringier_events",
	)
	if err != nil {
		log.Fatal(err)
	}
	handler := promhttp.Handler()
	go func(h http.Handler) {
		router := gin.Default()
		router.Any("/metrics", func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		})
		router.Run(":9000")
	}(handler)

	defer db.CloseDbConnection(ctx)
	log.Println("[persister] Starting gRPC server on port 9090")
	server, err := gRPC.NewService(
		gRPC.WithDatabase(db),
		gRPC.WithMetrics(&gRPC.Metrics{
			Latency:     persisterLatency,
			TotalEvents: persisterTotalReqs,
		}),
	)
	if err != nil {
		log.Fatalf("could not initialize service err: %s", err)
	}
	server.Run(":9090")
}
