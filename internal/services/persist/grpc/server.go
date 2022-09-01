package grpc

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weezqyd/ringier/internal/adapters/database"
	"github.com/weezqyd/ringier/internal/common/messages"
	"github.com/weezqyd/ringier/internal/services/persist/grpc/persister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Service struct {
	persister.UnimplementedPersistServiceServer

	db      *database.DB
	metrics *Metrics
}

type Metrics struct {
	Latency     *prometheus.HistogramVec
	TotalEvents *prometheus.CounterVec
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

func WithDatabase(db *database.DB) ServiceOption {
	return func(s *Service) error {
		s.db = db
		return nil
	}
}

func WithMetrics(metrics *Metrics) ServiceOption {
	return func(s *Service) error {
		if metrics.Latency == nil || metrics.TotalEvents == nil {
			return fmt.Errorf("missing metrics collectors")
		}
		s.metrics = metrics
		return nil
	}
}

// Run registers the Persist service to a grpcServer and serves on
// the specified port
func (s *Service) Run(port string) {
	var err error

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s %v", port, err)
	}
	sever := s
	grpcServer := grpc.NewServer()
	persister.RegisterPersistServiceServer(grpcServer, sever)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port %s %v", port, err)
	}
}

func (s *Service) Save(ctx context.Context, req *messages.EventRequest) (*messages.EventResponse, error) {
	// we can do validation here but for now we'll just save the event
	resp := &messages.EventResponse{
		RequestId: req.RequestId,
		Errors:    []string{},
	}
	var timer *prometheus.Timer
	if s.metrics != nil {
		timer = prometheus.NewTimer(s.metrics.Latency.WithLabelValues(req.Event))
		defer timer.ObserveDuration()
	}
	if err := s.db.SaveEvent(ctx, req.Event, req); err != nil {
		resp.Success = false
		resp.Errors = []string{fmt.Sprintf("failed to save event err: %s", err.Error())}
		if s.metrics != nil {
			s.metrics.TotalEvents.WithLabelValues(req.Event, "failed").Inc()
		}
		return resp, status.Error(codes.Internal, err.Error())
	}
	resp.Success = true
	if s.metrics != nil {
		s.metrics.TotalEvents.WithLabelValues(req.Event, "success").Inc()
	}

	return resp, nil
}
