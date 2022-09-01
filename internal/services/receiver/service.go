package receiver

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weezqyd/ringier/internal/adapters/msb"
	"github.com/weezqyd/ringier/internal/common"
	"log"
	"net/http"
	"strconv"
)

type Service struct {
	publisher *msb.Publisher
	config    *ServiceConfig
	metrics   *Metrics
}
type ServiceConfig struct {
	RabbitMQURI  string
	KafkaBrokers []string
	ListenPort   uint
	Topic        string
}

type Metrics struct {
	Latency       *prometheus.HistogramVec
	TotalRequests *prometheus.CounterVec
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type ServiceOption func(service *Service) error

// NewService creates a new receiver Service instance.
func NewService(opts ...ServiceOption) (*Service, error) {
	service := &Service{}
	for _, fn := range opts {
		if err := fn(service); err != nil {
			return nil, err
		}
	}

	return service, nil
}

func WithPublisher() ServiceOption {
	return func(service *Service) error {
		publisher, err := msb.NewPublisher(
			msb.WithRabbitPublisher(service.config.RabbitMQURI),
			msb.WithKafkaPublisher(service.config.KafkaBrokers),
			msb.WithTopic(service.config.Topic),
		)
		if err != nil {
			return err
		}
		service.publisher = publisher
		return nil
	}
}

func WithConfig(config *ServiceConfig) ServiceOption {
	return func(service *Service) error {
		service.config = config
		return nil
	}
}

func WithMetrics(metrics *Metrics) ServiceOption {
	return func(s *Service) error {
		if metrics.Latency == nil || metrics.TotalRequests == nil {
			return fmt.Errorf("missing metrics collectors")
		}
		s.metrics = metrics
		return nil
	}
}

func (r *Service) registerRoutes(app *gin.Engine) {
	app.Group("/api/v1").
		POST("/events", r.publishEvents)
}

func (r *Service) publishEvents(ctx *gin.Context) {
	request := common.EventRequest{}
	status := "failed"
	err := ctx.BindJSON(&request)
	if err != nil {
		respondError(ctx, []string{err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.metrics.TotalRequests.With(prometheus.Labels{"event": request.Event, "code": status}).Inc()
	timer := prometheus.NewTimer(r.metrics.Latency.With(
		prometheus.Labels{
			"event": request.Event,
		}),
	)
	defer timer.ObserveDuration()
	if err := validate.Struct(request); err != nil {
		var errors []string
		for _, msg := range err.(validator.ValidationErrors) {
			errors = append(errors, msg.Error())
		}
		respondError(ctx, errors, http.StatusUnprocessableEntity)
		return
	}
	event := common.Event{
		EventRequest: request,
		RequestId:    uuid.New().String(),
	}
	payload, err := json.Marshal(&event)
	if err != nil {
		respondError(ctx, []string{err.Error()}, http.StatusBadRequest)
		return
	}
	msg := buildMessage(payload)
	if err := r.publisher.Publish(msg); err != nil {
		respondError(ctx, []string{err.Error()}, http.StatusBadRequest)
		return
	}
	status = "success"
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"request_id": event.RequestId,
		"errors":     []string{},
	})
}

func respondError(ctx *gin.Context, errors []string, status int) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"success":    false,
		"request_id": "",
		"errors":     errors,
	})
	ctx.Status(status)
	ctx.Header("Content-Type", "application/json")
}

func (r *Service) Start(app *gin.Engine) {
	validate = validator.New()
	r.registerRoutes(app)
	log.Fatal(app.Run(":" + strconv.Itoa(int(r.config.ListenPort))))
}

func buildMessage(payload []byte) *message.Message {
	return message.NewMessage(
		uuid.New().String(),
		payload,
	)
}
