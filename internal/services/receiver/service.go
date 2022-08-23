package receiver

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/weezqyd/ringier/internal/adapters/msb"
	"github.com/weezqyd/ringier/internal/common"
	"log"
	"net/http"
	"strconv"
)

type Receiver struct {
	publisher *msb.Publisher
	config    *ServiceConfig
}
type ServiceConfig struct {
	RabbitMQURI  string
	KafkaBrokers []string
	ListenPort   uint
	Topic        string
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func NewReceiverService(config *ServiceConfig) (*Receiver, error) {
	receiver := &Receiver{}
	receiver.config = config
	publisher, err := msb.NewPublisher(
		msb.WithRabbitPublisher(config.RabbitMQURI),
		msb.WithKafkaPublisher(config.KafkaBrokers),
		msb.WithTopic(config.Topic),
	)
	if err != nil {
		return nil, err
	}
	receiver.publisher = publisher

	return receiver, nil
}

func (r Receiver) registerRoutes(app *gin.Engine) {
	app.Group("/api/v1").
		POST("/events", r.publishEvents)
}

func (r Receiver) publishEvents(ctx *gin.Context) {
	request := common.EventRequest{}
	body, err := ctx.GetRawData()
	if err != nil {
		respondError(ctx, []string{err.Error()}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &request); err != nil {
		respondError(ctx, []string{err.Error()}, http.StatusBadRequest)
		return
	}
	log.Println(request)

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

func (r Receiver) Start(app *gin.Engine) {
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
