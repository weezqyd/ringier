package grpc

import (
	"context"
	"fmt"
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

	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
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
	if err := s.db.SaveEvent(ctx, req.Event, req); err != nil {
		resp.Success = false
		resp.Errors = []string{fmt.Sprintf("failed to save event err: %s", err.Error())}

		return resp, status.Error(codes.Internal, err.Error())
	}
	resp.Success = true

	return resp, nil
}
